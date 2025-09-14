import {Button} from "@/components/ui/button";
import {Card, CardHeader, CardTitle} from "@/components/ui/card";
import sessionService from "@/services/session.service";
import {useSessionStore} from "@/stores/session.store";
import type {Session} from "@/types/models/session";
import {Loader2, Maximize, Video, VideoOff} from "lucide-react";
import React, {useEffect, useRef, useState} from "react";
import {useNavigate, useParams, useSearchParams} from "react-router-dom";

const SessionPage: React.FC = () => {
  const {sessionId} = useParams();
  const [query] = useSearchParams();

  const navigate = useNavigate();

  const {currentSession, setCurrentSession, startSession, host} =
    useSessionStore();

  const hostId = query.get("hostId") || "";
  const clientId = query.get("clientId") || "";

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [starting, setStarting] = useState(false);
  const [connected, setConnected] = useState(false);
  const [videoEnabled, setVideoEnabled] = useState(true);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);
  const [statusMessage, setStatusMessage] = useState("Connecting...");

  // Add session ending state
  const [isEndingSession, setIsEndingSession] = useState(false);

  // Add reconnection state
  const [reconnectAttempts, setReconnectAttempts] = useState(0);
  const [maxReconnectAttempts] = useState(5);
  const [reconnectTimeout, setReconnectTimeout] =
    useState<NodeJS.Timeout | null>(null);
  const [isReconnecting, setIsReconnecting] = useState(false);

  const pcRef = useRef<RTCPeerConnection | null>(null);
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [videoStream, setVideoStream] = useState<MediaStream | null>(null);

  // WebSocket connection with reconnection logic
  const connectWebSocket = (isReconnect = false) => {
    if (
      !currentSession ||
      !host?.ip_address ||
      !host?.port ||
      isEndingSession
    ) {
      console.log(
        "Cannot connect WebSocket: missing session, host info, or session is ending"
      );
      return;
    }

    // Don't reconnect if we've exceeded max attempts or if session is ending
    if (
      isReconnect &&
      (reconnectAttempts >= maxReconnectAttempts || isEndingSession)
    ) {
      console.log("Max reconnection attempts reached or session ending");
      setStatusMessage("Connection failed - max retries exceeded");
      setIsReconnecting(false);
      return;
    }

    if (isReconnect) {
      setIsReconnecting(true);
      setStatusMessage(
        `Reconnecting... (${reconnectAttempts + 1}/${maxReconnectAttempts})`
      );
    }

    const wsUrl = `ws://${host.ip_address}:${host.port}/ws?session_id=${currentSession.id}`;
    console.log("Connecting to WebSocket:", wsUrl);

    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
      setWsConnected(true);
      setIsReconnecting(false);
      setReconnectAttempts(0); // Reset attempts on successful connection
      setStatusMessage("Connected - Ready for input");

      // Clear any pending reconnection timeout
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        setReconnectTimeout(null);
      }
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("WebSocket message:", data);

      if (data.type === "status") {
        setStatusMessage(data.message || "Connected");
      }
    };

    ws.onclose = (event) => {
      console.log("WebSocket disconnected", event);
      setWsConnected(false);
      setIsReconnecting(false);

      // Don't attempt reconnection if it was a clean close or if we're ending the session
      if (event.wasClean || !currentSession) {
        setStatusMessage("Disconnected");
        return;
      }

      // Attempt reconnection
      if (reconnectAttempts < maxReconnectAttempts) {
        const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000); // Exponential backoff, max 30s
        console.log(
          `WebSocket disconnected, attempting reconnection in ${delay}ms (attempt ${
            reconnectAttempts + 1
          }/${maxReconnectAttempts})`
        );

        setStatusMessage(
          `Connection lost, reconnecting in ${Math.round(delay / 1000)}s...`
        );

        const timeout = setTimeout(() => {
          setReconnectAttempts((prev) => prev + 1);
          connectWebSocket(true);
        }, delay);

        setReconnectTimeout(timeout);
      } else {
        setStatusMessage("Connection lost - reconnection failed");
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      if (!isReconnecting) {
        setStatusMessage("Connection error");
      }
    };
  };

  // Manual reconnection function
  const reconnectWebSocket = () => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    setReconnectAttempts(0);
    setIsReconnecting(false);
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      setReconnectTimeout(null);
    }
    connectWebSocket();
  };

  // Cleanup function
  const cleanupWebSocket = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      setReconnectTimeout(null);
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setWsConnected(false);
    setIsReconnecting(false);
    setReconnectAttempts(0);
  };

  // Send input events to host
  const sendInputEvent = (command: any) => {
    if (
      wsRef.current &&
      wsRef.current.readyState === WebSocket.OPEN &&
      currentSession
    ) {
      const message = {
        sessionId: currentSession.id,
        command: command,
      };
      console.log("Sending input event:", message);
      wsRef.current.send(JSON.stringify(message));
    }
  };

  // Mouse event handlers
  const handleMouseMove = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;

    const rect = videoRef.current.getBoundingClientRect();
    const x = Math.round((e.clientX - rect.left) * (rect.width / rect.width));
    const y = Math.round((e.clientY - rect.top) * (rect.height / rect.height));

    sendInputEvent({
      type: "mouse",
      action: "move",
      x: x,
      y: y,
    });
  };

  const handleMouseClick = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;

    const rect = videoRef.current.getBoundingClientRect();
    const x = Math.round((e.clientX - rect.left) * (rect.width / rect.width));
    const y = Math.round((e.clientY - rect.top) * (rect.height / rect.height));

    const buttonMap = {0: "left", 1: "middle", 2: "right"};
    const button = buttonMap[e.button as keyof typeof buttonMap] || "left";

    sendInputEvent({
      type: "mouse",
      action: "click",
      button: button,
      x: x,
      y: y,
    });
  };

  const handleMouseDown = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;

    const rect = videoRef.current.getBoundingClientRect();
    const x = Math.round((e.clientX - rect.left) * (rect.width / rect.width));
    const y = Math.round((e.clientY - rect.top) * (rect.height / rect.height));

    const buttonMap = {0: "left", 1: "middle", 2: "right"};
    const button = buttonMap[e.button as keyof typeof buttonMap] || "left";

    sendInputEvent({
      type: "mouse",
      action: "press",
      button: button,
      x: x,
      y: y,
    });
  };

  const handleMouseUp = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;

    const rect = videoRef.current.getBoundingClientRect();
    const x = Math.round((e.clientX - rect.left) * (rect.width / rect.width));
    const y = Math.round((e.clientY - rect.top) * (rect.height / rect.height));

    const buttonMap = {0: "left", 1: "middle", 2: "right"};
    const button = buttonMap[e.button as keyof typeof buttonMap] || "left";

    sendInputEvent({
      type: "mouse",
      action: "release",
      button: button,
      x: x,
      y: y,
    });
  };

  // Keyboard event handlers
  const handleKeyDown = (e: React.KeyboardEvent<HTMLVideoElement>) => {
    e.preventDefault();

    sendInputEvent({
      type: "keyboard",
      action: "press",
      key: e.key,
    });
  };

  const handleKeyUp = (e: React.KeyboardEvent<HTMLVideoElement>) => {
    e.preventDefault();

    sendInputEvent({
      type: "keyboard",
      action: "release",
      key: e.key,
    });
  };

  // Fullscreen handling
  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      containerRef.current?.requestFullscreen();
      setIsFullscreen(true);
    } else {
      document.exitFullscreen();
      setIsFullscreen(false);
    }
  };

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener("fullscreenchange", handleFullscreenChange);
    return () =>
      document.removeEventListener("fullscreenchange", handleFullscreenChange);
  }, []);

  useEffect(() => {
    if (!sessionId) {
      navigate("/hosts", {replace: true});
      return;
    }

    let canceled = false;

    const run = async () => {
      setLoading(true);
      setError(null);
      setStatusMessage("Loading session...");

      try {
        const s = await sessionService.get(sessionId);
        if (canceled) return;
        setCurrentSession(s);

        if (s.status === "ACTIVE") {
          setStatusMessage("Connecting to video...");
          await establishWebRTCConnection(s);
          return;
        }

        if (s.status === "PENDING") {
          setStarting(true);
          setStatusMessage("Starting session...");

          const pc = new RTCPeerConnection({
            iceServers: [{urls: "stun:stun.l.google.com:19302"}],
          });
          pcRef.current = pc;

          setupWebRTCEventHandlers(pc);
          pc.addTransceiver("video", {direction: "recvonly"});

          const offer = await pc.createOffer();
          await pc.setLocalDescription(offer);

          if (!offer.sdp) {
            throw new Error("Failed to create offer");
          }

          setStatusMessage("Connecting to host...");
          await startSession({webrtc_offer: offer.sdp}, sessionId);

          const activeSession = await sessionService.pollUntilActive(
            sessionId,
            (updated) => {
              if (canceled) return;
              setCurrentSession(updated);
            }
          );

          setStatusMessage("Establishing video connection...");
          await establishWebRTCConnection(activeSession);
        }
      } catch (e: unknown) {
        setError(e instanceof Error ? e.message : "Unknown error");
        setStatusMessage("Error occurred");
      } finally {
        setStarting(false);
        setLoading(false);
      }
    };
    run();

    return () => {
      canceled = true;
      pcRef.current?.close();
      pcRef.current = null;
      cleanupWebSocket();
    };
  }, [sessionId, navigate, startSession, setCurrentSession]);

  const setupWebRTCEventHandlers = (pc: RTCPeerConnection) => {
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        console.log("ICE candidate:", event.candidate);
      }
    };

    pc.ontrack = (event) => {
      console.log("Received remote track:", event);

      if (event.streams && event.streams[0]) {
        const stream = event.streams[0];
        setVideoStream(stream);
        setConnected(true);
        setStatusMessage("Video connected");
      }
    };

    pc.onconnectionstatechange = () => {
      console.log("Connection state:", pc.connectionState);
      if (pc.connectionState === "connected") {
        setConnected(true);
        setStatusMessage("Video connected");
      } else if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        setConnected(false);
        setStatusMessage("Video disconnected");
      }
    };

    pc.oniceconnectionstatechange = () => {
      console.log("ICE connection state:", pc.iceConnectionState);
    };

    pc.onsignalingstatechange = () => {
      console.log("Signaling state:", pc.signalingState);
    };
  };

  const establishWebRTCConnection = async (session: Session) => {
    if (!session.webrtc_answer) {
      throw new Error("No WebRTC answer received from host");
    }

    const pc = pcRef.current;
    if (!pc) {
      throw new Error("PeerConnection not initialized");
    }

    try {
      await pc.setRemoteDescription({
        type: "answer",
        sdp: session.webrtc_answer,
      });

      console.log("WebRTC connection established");
    } catch (error) {
      console.error("Failed to establish WebRTC connection:", error);
      throw error;
    }
  };

  const toggleVideo = () => {
    setVideoEnabled(!videoEnabled);
  };

  const endSession = async () => {
    try {
      setIsEndingSession(true);
      setStatusMessage("Ending session...");

      // Clean up WebSocket connection first
      cleanupWebSocket();

      // End the session
      await sessionService.end({reason: "Ended by user"}, sessionId!);

      // Navigate away
      navigate("/hosts");
    } catch (error) {
      console.error("Failed to end session:", error);
      setError("Failed to end session");
      setIsEndingSession(false); // Reset flag on error
    }
  };

  useEffect(() => {
    if (videoStream && videoRef.current) {
      videoRef.current.srcObject = videoStream;
      videoRef.current
        .play()
        .then(() => {
          console.log("Video started playing, connecting WebSocket...");
          connectWebSocket();
        })
        .catch((err) => {
          console.error("Failed to play video:", err);
        });
    }
  }, [videoStream]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      cleanupWebSocket();
    };
  }, []);

  return (
    <div
      ref={containerRef}
      className={`${isFullscreen ? "fixed inset-0 bg-black z-50" : ""}`}
    >
      {/* Status Bar */}
      <div
        className={`${
          isFullscreen
            ? "fixed bottom-0 left-0 right-0 bg-black/80 text-white p-2 z-50"
            : "mb-4"
        }`}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4 text-sm">
            <div className="flex items-center gap-1">
              <div
                className={`w-2 h-2 rounded-full ${
                  connected ? "bg-green-500" : "bg-red-500"
                }`}
              />
              <span>Video: {connected ? "Connected" : "Disconnected"}</span>
            </div>
            <div className="flex items-center gap-1">
              <div
                className={`w-2 h-2 rounded-full ${
                  wsConnected
                    ? "bg-green-500"
                    : isReconnecting
                    ? "bg-yellow-500 animate-pulse"
                    : "bg-red-500"
                }`}
              />
              <span>
                Input:{" "}
                {wsConnected
                  ? "Connected"
                  : isReconnecting
                  ? "Reconnecting..."
                  : "Disconnected"}
              </span>
            </div>
            <span>{statusMessage}</span>
          </div>

          {/* Reconnection controls - hide when ending session */}
          {!isEndingSession &&
            !wsConnected &&
            !isReconnecting &&
            reconnectAttempts < maxReconnectAttempts && (
              <button
                onClick={reconnectWebSocket}
                className="px-3 py-1 bg-blue-600 text-white text-xs rounded hover:bg-blue-700"
              >
                Reconnect
              </button>
            )}

          {!isEndingSession &&
            !wsConnected &&
            reconnectAttempts >= maxReconnectAttempts && (
              <button
                onClick={reconnectWebSocket}
                className="px-3 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700"
              >
                Retry Connection
              </button>
            )}

          {/* End Session button */}
          {!isEndingSession && (
            <button
              onClick={endSession}
              className="px-3 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700"
            >
              End Session
            </button>
          )}

          {/* Fullscreen toggle */}
          <button
            onClick={toggleFullscreen}
            className="px-3 py-1 bg-gray-600 text-white text-xs rounded hover:bg-gray-700"
          >
            {isFullscreen ? "Exit Fullscreen" : "Fullscreen"}
          </button>
        </div>
      </div>

      {/* Main Content */}
      {!isFullscreen && (
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span>Session {sessionId}</span>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={toggleVideo}
                  disabled={!connected}
                >
                  {videoEnabled ? (
                    <VideoOff className="w-4 h-4" />
                  ) : (
                    <Video className="w-4 h-4" />
                  )}
                </Button>
                <Button variant="outline" size="sm" onClick={toggleFullscreen}>
                  <Maximize className="w-4 h-4" />
                </Button>
              </div>
            </CardTitle>
          </CardHeader>
        </Card>
      )}

      {(loading || starting) && (
        <div className="flex justify-center items-center h-64">
          <Loader2 className="w-8 h-8 animate-spin" />
          <span className="ml-2 text-lg">{statusMessage}</span>
        </div>
      )}

      {error && (
        <div className="text-destructive p-4 border border-destructive rounded mb-4">
          {error}
        </div>
      )}

      {!loading &&
        !error &&
        currentSession &&
        currentSession.status === "ACTIVE" && (
          <div className={`${isFullscreen ? "h-full" : ""}`}>
            {connected ? (
              <video
                ref={videoRef}
                autoPlay
                playsInline
                muted
                className={`${
                  isFullscreen
                    ? "w-full h-full object-contain"
                    : "w-full max-w-4xl mx-auto rounded-lg border"
                }`}
                style={{display: videoEnabled ? "block" : "none"}}
                onMouseMove={handleMouseMove}
                onClick={handleMouseClick}
                onMouseDown={handleMouseDown}
                onMouseUp={handleMouseUp}
                onKeyDown={handleKeyDown}
                onKeyUp={handleKeyUp}
                tabIndex={0}
              />
            ) : (
              <div className="flex items-center justify-center h-64 border-2 border-dashed border-muted-foreground rounded-lg">
                <div className="text-center">
                  <VideoOff className="w-12 h-12 mx-auto mb-2 text-muted-foreground" />
                  <p className="text-muted-foreground">
                    {starting
                      ? "Connecting to video feed..."
                      : "No video feed available"}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

      {/* Session Info (only in windowed mode) */}
      {!isFullscreen && !loading && !error && currentSession && (
        <div className="grid grid-cols-2 gap-4 text-sm text-muted-foreground mt-4">
          <div>
            <strong>Status:</strong> {currentSession.status}
          </div>
          <div>
            <strong>Host:</strong> {currentSession.host_name}
          </div>
          <div>
            <strong>Client:</strong> {currentSession.client_name}
          </div>
          <div>
            <strong>Session ID:</strong> {sessionId}
          </div>
        </div>
      )}
    </div>
  );
};

export default SessionPage;
