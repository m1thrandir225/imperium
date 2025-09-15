import {Button} from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import sessionService from "@/services/session.service";
import {useSessionStore} from "@/stores/session.store";
import type {Session} from "@/types/models/session";
import {Loader2} from "lucide-react";
import React, {useEffect, useRef, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";

const SessionPage: React.FC = () => {
  const {sessionId} = useParams();

  const navigate = useNavigate();

  const {currentSession, setCurrentSession, startSession} = useSessionStore();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [starting, setStarting] = useState(false);
  const [connected, setConnected] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [statusMessage, setStatusMessage] = useState("Connecting...");
  const [isEndingSession, setIsEndingSession] = useState(false);

  const pcRef = useRef<RTCPeerConnection | null>(null);
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [videoStream, setVideoStream] = useState<MediaStream | null>(null);
  const [audioStream, setAudioStream] = useState<MediaStream | null>(null);

  // DataChannel for input
  const dataChannelRef = useRef<RTCDataChannel | null>(null);
  const [dcConnected, setDcConnected] = useState(false);

  const toU16 = (n: number) =>
    Math.max(0, Math.min(65535, Math.round(n * 65535)));
  const mapMouseButton = (btn: number) => (btn === 1 ? 3 : btn === 2 ? 2 : 1);

  const VK: Record<string, number> = {
    Escape: 0x1b,
    Enter: 0x0d,
    Space: 0x20,
    ShiftLeft: 0x10,
    ShiftRight: 0x10,
    ControlLeft: 0x11,
    ControlRight: 0x11,
    AltLeft: 0x12,
    AltRight: 0x12,
    ArrowUp: 0x26,
    ArrowDown: 0x28,
    ArrowLeft: 0x25,
    ArrowRight: 0x27,
    KeyW: 0x57,
    KeyA: 0x41,
    KeyS: 0x53,
    KeyD: 0x44,
    KeyQ: 0x51,
    KeyE: 0x45,
    KeyR: 0x52,
    KeyF: 0x46,
    KeyG: 0x47,
    KeyH: 0x48,
    KeyI: 0x49,
    KeyJ: 0x4a,
    KeyK: 0x4b,
    KeyL: 0x4c,
    KeyM: 0x4d,
    KeyN: 0x4e,
    KeyO: 0x4f,
    KeyP: 0x50,
    KeyZ: 0x5a,
    KeyX: 0x58,
    KeyY: 0x59,
    Digit1: 0x31,
    Digit2: 0x32,
    Digit3: 0x33,
    Digit4: 0x34,
    Digit5: 0x35,
    Digit6: 0x36,
    Digit7: 0x37,
    Digit8: 0x38,
    Digit9: 0x39,
    Digit0: 0x30,
    Semicolon: 0x3a,
    Equals: 0x3b,
    Comma: 0x3c,
    Minus: 0x3d,
    Period: 0x3e,
    Slash: 0x3f,
    Backquote: 0x60,
    Openbracket: 0x5b,
    Backslash: 0x5c,
    Closebracket: 0x5d,
    Caret: 0x5e,
    Underscore: 0x5f,
  };
  const getVK = (e: React.KeyboardEvent) => VK[e.code] ?? 0;

  const sendFrame10 = (
    t: number,
    a: number,
    btn: number,
    vk: number,
    x: number,
    y: number
  ) => {
    const dc = dataChannelRef.current;
    if (!dc || dc.readyState !== "open") return;
    const buf = new ArrayBuffer(10);
    const dv = new DataView(buf);
    dv.setUint8(0, t);
    dv.setUint8(1, a);
    dv.setUint8(2, btn);
    dv.setUint8(3, 0);
    dv.setUint16(4, vk, true);
    dv.setUint16(6, x, true);
    dv.setUint16(8, y, true);
    dc.send(buf);
  };

  const rel = (e: React.MouseEvent<HTMLVideoElement>) => {
    const v = videoRef.current!;
    const r = v.getBoundingClientRect();

    const vw = v.videoWidth || 1;
    const vh = v.videoHeight || 1;

    // Element and video aspect ratios
    const arElem = r.width / r.height;
    const arVid = vw / vh;

    // Compute rendered content rect inside the element (object-contain)
    let contentW: number,
      contentH: number,
      offsetX = 0,
      offsetY = 0;
    if (arElem > arVid) {
      // pillarbox: height fits, width centered
      contentH = r.height;
      contentW = contentH * arVid;
      offsetX = (r.width - contentW) / 2;
    } else {
      // letterbox: width fits, height centered
      contentW = r.width;
      contentH = contentW / arVid;
      offsetY = (r.height - contentH) / 2;
    }

    // Mouse position relative to content box
    const px = e.clientX - r.left - offsetX;
    const py = e.clientY - r.top - offsetY;

    // Normalize to 0..1, clamp
    const x = Math.max(0, Math.min(1, px / contentW));
    const y = Math.max(0, Math.min(1, py / contentH));

    // Debug (optional)
    // console.log("[REL]", {px, py, contentW, contentH, offsetX, offsetY, x, y});

    return {x, y};
  };

  const sendMouseMove = (xNorm: number, yNorm: number) =>
    sendFrame10(1, 2, 0, 0, toU16(xNorm), toU16(yNorm));
  const sendMouseDown = (button: number, xNorm: number, yNorm: number) =>
    sendFrame10(2, 0, mapMouseButton(button), 0, toU16(xNorm), toU16(yNorm));
  const sendMouseUp = (button: number, xNorm: number, yNorm: number) =>
    sendFrame10(2, 1, mapMouseButton(button), 0, toU16(xNorm), toU16(yNorm));
  const sendMouseClick = (button: number, xNorm: number, yNorm: number) =>
    sendFrame10(2, 3, mapMouseButton(button), 0, toU16(xNorm), toU16(yNorm));
  const sendWheel = (deltaY: number) => {
    const dc = dataChannelRef.current;
    if (!dc || dc.readyState !== "open") return;
    const buf = new ArrayBuffer(10);
    const dv = new DataView(buf);
    dv.setUint8(0, 3); // wheel
    dv.setUint8(1, 0);
    dv.setUint8(2, 0);
    dv.setUint8(3, 0);
    dv.setUint16(4, 0, true);
    dv.setUint16(6, 0, true);
    dv.setInt16(8, Math.max(-32768, Math.min(32767, deltaY)), true);
    dc.send(buf);
  };

  // ---- Input handlers (video element) ----
  const handleMouseMove = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;
    const {x, y} = rel(e);
    sendMouseMove(x, y);
  };
  const handleMouseDown = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;
    const {x, y} = rel(e);
    sendMouseDown(e.button, x, y);
  };
  const handleMouseUp = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;
    const {x, y} = rel(e);
    sendMouseUp(e.button, x, y);
  };
  const handleMouseClick = (e: React.MouseEvent<HTMLVideoElement>) => {
    if (!videoRef.current) return;
    const {x, y} = rel(e);
    sendMouseClick(e.button, x, y);
  };
  const handleWheel = (e: React.WheelEvent<HTMLVideoElement>) => {
    sendWheel(e.deltaY);
  };
  const handleKeyDown = (e: React.KeyboardEvent<HTMLVideoElement>) => {
    e.preventDefault();
    sendFrame10(0, 0, 0, getVK(e), 0, 0);
  };
  const handleKeyUp = (e: React.KeyboardEvent<HTMLVideoElement>) => {
    e.preventDefault();
    sendFrame10(0, 1, 0, getVK(e), 0, 0);
  };

  // ---- WebRTC setup ----
  const setupWebRTCEventHandlers = (pc: RTCPeerConnection) => {
    pc.onicecandidate = (event) => {
      if (event.candidate) console.log("ICE candidate:", event.candidate);
    };

    pc.ondatachannel = (e) => {
      const dc = e.channel;
      if (dc.label !== "input") {
        console.log("Ignoring unexpected DataChannel:", dc.label);
        return;
      }

      dataChannelRef.current = dc;

      dc.onopen = () => {
        setDcConnected(true);
        console.log("[DC] open", {
          label: dc.label,
          id: dc.id,
          state: dc.readyState,
        });
      };
      dc.onclose = () => {
        setDcConnected(false);
        console.log("[DC] close", {label: dc.label, id: dc.id});
      };
      dc.onerror = (err) => {
        setDcConnected(false);
        console.log("[DC] error", err);
      };
      dc.onmessage = (msg) => {
        // Optional: handle acks/status if host sends any
        console.log(
          "[DC] message",
          msg.data instanceof ArrayBuffer ? new Uint8Array(msg.data) : msg.data
        );
      };
    };

    pc.ontrack = (event) => {
      if (event.streams && event.streams[0]) {
        const stream = event.streams[0];
        setVideoStream(stream);
        setConnected(true);
        setStatusMessage("Video connected");
        if (event.track.kind === "audio") setAudioStream(stream);
      }
    };

    pc.onconnectionstatechange = () => {
      if (pc.connectionState === "connected") {
        setConnected(true);
        setStatusMessage("Connected");
      } else if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        setConnected(false);
        setStatusMessage("Disconnected");
        setDcConnected(false);
      }
    };
  };

  const establishWebRTCConnection = async (session: Session) => {
    if (!session.webrtc_answer)
      throw new Error("No WebRTC answer received from host");
    const pc = pcRef.current;
    if (!pc) throw new Error("PeerConnection not initialized");
    await pc.setRemoteDescription({type: "answer", sdp: session.webrtc_answer});
    console.log("WebRTC connection established");
  };

  const endSession = async () => {
    try {
      setIsEndingSession(true);
      setStatusMessage("Ending session...");

      try {
        dataChannelRef.current?.close();
      } catch {
        console.error("Failed to close data channel");
      }
      dataChannelRef.current = null;

      try {
        await pcRef.current?.close();
      } catch {
        console.error("Failed to close peer connection");
      }
      pcRef.current = null;

      await sessionService.end({reason: "Ended by user"}, sessionId!);
      navigate("/hosts");
    } catch (err) {
      console.error("Failed to end session:", err);
      setError("Failed to end session");
      setIsEndingSession(false);
    }
  };

  // ---- Media hookup ----
  useEffect(() => {
    if (videoStream && videoRef.current) {
      videoRef.current.srcObject = videoStream;
      videoRef.current
        .play()
        .catch((err) => console.error("Failed to play video:", err));
    }
  }, [videoStream]);

  useEffect(() => {
    if (audioStream && audioRef.current) {
      audioRef.current.srcObject = audioStream;
      audioRef.current
        .play()
        .catch((err) => console.error("Failed to play audio:", err));
    }
  }, [audioStream]);

  // ---- Fullscreen ----
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
    const onFs = () => setIsFullscreen(!!document.fullscreenElement);
    document.addEventListener("fullscreenchange", onFs);
    return () => document.removeEventListener("fullscreenchange", onFs);
  }, []);

  // ---- Session lifecycle ----
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
          const pc = new RTCPeerConnection({
            iceServers: [{urls: "stun:stun.l.google.com:19302"}],
          });
          pcRef.current = pc;
          setupWebRTCEventHandlers(pc);
          // When resuming an already-active session, we still need to set remote answer
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
          pc.addTransceiver("audio", {direction: "recvonly"});

          // Ensure SCTP is negotiated in the offer (host creates "input")
          const bootstrap = pc.createDataChannel("_bootstrap"); // negotiated: false (default)
          bootstrap.onopen = () => {
            console.log("[DC] bootstrap open → closing");
            try {
              bootstrap.close();
            } catch {
              console.error("Failed to close bootstrap data channel");
            }
          };
          bootstrap.onclose = () => console.log("[DC] bootstrap closed");

          // Now create the offer
          const offer = await pc.createOffer();
          await pc.setLocalDescription(offer);
          if (!offer.sdp) throw new Error("Failed to create offer");

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
        setLoading(false);
        setStarting(false);
      }
    };

    run();

    return () => {
      canceled = true;
    };
  }, [sessionId, navigate, startSession, setCurrentSession]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      try {
        dataChannelRef.current?.close();
      } catch {
        console.error("Failed to close data channel");
      }
      dataChannelRef.current = null;
      try {
        pcRef.current?.close();
      } catch {
        console.error("Failed to close peer connection");
      }
      pcRef.current = null;
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
                  audioStream ? "bg-green-500" : "bg-red-500"
                }`}
              />
              <span>Audio: {audioStream ? "Connected" : "Disconnected"}</span>
            </div>
            <div className="flex items-center gap-1">
              <div
                className={`w-2 h-2 rounded-full ${
                  dcConnected ? "bg-green-500" : "bg-red-500"
                }`}
              />
              <span>Input: {dcConnected ? "Connected" : "Disconnected"}</span>
            </div>
            <span>{statusMessage}</span>
          </div>

          <div className="flex items-center gap-2">
            {!isEndingSession && (
              <button
                onClick={endSession}
                className="px-3 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700"
              >
                End Session
              </button>
            )}
            <button
              onClick={toggleFullscreen}
              className="px-3 py-1 bg-gray-600 text-white text-xs rounded hover:bg-gray-700"
            >
              {isFullscreen ? "Exit Fullscreen" : "Fullscreen"}
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      {!isFullscreen && (
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span>Session {sessionId}</span>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={toggleFullscreen}>
                  Fullscreen
                </Button>
              </div>
            </CardTitle>
            <CardDescription>
              Host: {currentSession?.host_name} | Status:{" "}
              {currentSession?.status}
            </CardDescription>
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
              <div className="relative">
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
                  style={{display: "block"}}
                  onMouseMove={handleMouseMove}
                  onClick={handleMouseClick}
                  onMouseDown={handleMouseDown}
                  onMouseUp={handleMouseUp}
                  onWheel={handleWheel}
                  onKeyDown={handleKeyDown}
                  onKeyUp={handleKeyUp}
                  tabIndex={0}
                />
                <audio ref={audioRef} autoPlay style={{display: "none"}} />
              </div>
            ) : (
              <div className="flex items-center justify-center h-64 border-2 border-dashed border-muted-foreground rounded-lg">
                <div className="text-center">
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
