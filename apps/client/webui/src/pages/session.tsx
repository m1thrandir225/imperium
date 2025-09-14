import {Card, CardHeader, CardTitle} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import sessionService from "@/services/session.service";
import {useSessionStore} from "@/stores/session.store";
import {Loader2, Video, VideoOff} from "lucide-react";
import React from "react";
import {useEffect, useRef, useState} from "react";
import {useNavigate, useParams, useSearchParams} from "react-router-dom";
import type {Session} from "@/types/models/session";

const SessionPage: React.FC = () => {
  const {sessionId} = useParams();
  const [query] = useSearchParams();

  const navigate = useNavigate();

  const {currentSession, setCurrentSession, startSession} = useSessionStore();

  const hostId = query.get("hostId") || "";
  const clientId = query.get("clientId") || "";

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [starting, setStarting] = useState(false);
  const [connected, setConnected] = useState(false);
  const [videoEnabled, setVideoEnabled] = useState(true);

  const pcRef = useRef<RTCPeerConnection | null>(null);
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const [videoStream, setVideoStream] = useState<MediaStream | null>(null);

  useEffect(() => {
    if (!sessionId) {
      navigate("/hosts", {replace: true});
      return;
    }

    let canceled = false;

    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const s = await sessionService.get(sessionId);
        if (canceled) return;
        setCurrentSession(s);

        if (s.status === "ACTIVE") {
          // Session is already active, establish WebRTC connection
          await establishWebRTCConnection(s);
          return;
        }

        if (s.status === "PENDING") {
          setStarting(true);
          const pc = new RTCPeerConnection({
            iceServers: [{urls: "stun:stun.l.google.com:19302"}],
          });
          pcRef.current = pc;

          // Set up event handlers
          setupWebRTCEventHandlers(pc);

          // Add transceiver for receiving video
          pc.addTransceiver("video", {direction: "recvonly"});

          const offer = await pc.createOffer();
          await pc.setLocalDescription(offer);

          if (!offer.sdp) {
            throw new Error("Failed to create offer");
          }

          await startSession({webrtc_offer: offer.sdp}, sessionId);

          // Poll until session becomes active
          const activeSession = await sessionService.pollUntilActive(
            sessionId,
            (updated) => {
              if (canceled) return;
              setCurrentSession(updated);
            }
          );

          // Establish WebRTC connection with the answer
          await establishWebRTCConnection(activeSession);
        }
      } catch (e: unknown) {
        setError(e instanceof Error ? e.message : "Unknown error");
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
      console.log("Track kind:", event.track.kind);
      console.log("Track enabled:", event.track.enabled);
      console.log("Track readyState:", event.track.readyState);
      console.log("Streams:", event.streams);

      if (event.streams && event.streams[0]) {
        const stream = event.streams[0];
        console.log("Stream tracks:", stream.getTracks());
        console.log("Stream active:", stream.active);

        const videoTracks = stream.getVideoTracks();
        console.log("Video tracks:", videoTracks);

        if (videoTracks.length > 0) {
          console.log("Video track settings:", videoTracks[0].getSettings());
          console.log(
            "Video track constraints:",
            videoTracks[0].getConstraints()
          );
        }

        setVideoStream(stream);
        setConnected(true);
      }
    };

    pc.onconnectionstatechange = () => {
      console.log("Connection state:", pc.connectionState);
      if (pc.connectionState === "connected") {
        setConnected(true);
      } else if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        setConnected(false);
      }
    };

    // Add more debugging
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
      // Set the remote description with the answer from the host
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
    if (videoRef.current) {
      videoRef.current.style.display = videoEnabled ? "none" : "block";
    }
  };

  const endSession = async () => {
    try {
      await sessionService.end({reason: "Ended by user"}, sessionId!);
      navigate("/hosts");
    } catch (error) {
      console.error("Failed to end session:", error);
      setError("Failed to end session");
    }
  };

  useEffect(() => {
    if (videoStream && videoRef.current) {
      console.log("Setting video srcObject");
      videoRef.current.srcObject = videoStream;

      videoRef.current.onloadedmetadata = () => {
        console.log("Video metadata loaded");
        console.log(
          "Video dimensions:",
          videoRef.current?.videoWidth,
          "x",
          videoRef.current?.videoHeight
        );
        console.log("Video duration:", videoRef.current?.duration);
      };

      videoRef.current.oncanplay = () => {
        console.log("Video can play");
      };

      videoRef.current.onplay = () => {
        console.log("Video started playing");
      };

      videoRef.current.onpause = () => {
        console.log("Video paused");
      };

      videoRef.current.onended = () => {
        console.log("Video ended");
      };

      videoRef.current.onstalled = () => {
        console.log("Video stalled");
      };

      videoRef.current.onwaiting = () => {
        console.log("Video waiting");
      };

      videoRef.current.onerror = (e) => {
        console.error("Video error:", e);
        console.error("Video error details:", videoRef.current?.error);
      };

      // Monitor video track state
      const videoTracks = videoStream.getVideoTracks();
      if (videoTracks.length > 0) {
        const track = videoTracks[0];
        console.log("Monitoring video track:", track.id);

        track.onended = () => {
          console.log("Video track ended");
        };

        track.onmute = () => {
          console.log("Video track muted");
        };

        track.onunmute = () => {
          console.log("Video track unmuted");
        };
      }

      videoRef.current.play().catch((err) => {
        console.error("Failed to play video:", err);
      });
    }
  }, [videoStream]);

  return (
    <React.Fragment>
      <Card>
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
              <Button
                variant="destructive"
                size="sm"
                onClick={endSession}
                disabled={!currentSession || currentSession.status !== "ACTIVE"}
              >
                End Session
              </Button>
            </div>
          </CardTitle>
        </CardHeader>
      </Card>

      {(loading || starting) && (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="w-4 h-4 animate-spin" />
          <span className="ml-2">
            {starting ? "Starting session..." : "Loading..."}
          </span>
        </div>
      )}

      {error && (
        <div className="text-destructive p-4 border border-destructive rounded">
          {error}
        </div>
      )}

      {!loading && !error && currentSession && (
        <div className="grid gap-4">
          {/* Session Info */}
          <div className="grid grid-cols-2 gap-4 text-sm text-muted-foreground">
            <div>
              <strong>Status:</strong> {currentSession.status}
            </div>
            <div>
              <strong>Host:</strong> {currentSession.host_name} ({hostId})
            </div>
            <div>
              <strong>Client:</strong> {currentSession.client_name} ({clientId})
            </div>
            <div>
              <strong>Connection:</strong>{" "}
              {connected ? "Connected" : "Disconnected"}
            </div>
          </div>

          {/* Video Feed */}
          {currentSession.status === "ACTIVE" && (
            <Card>
              <CardHeader>
                <CardTitle>Video Feed</CardTitle>
              </CardHeader>
              <div className="p-4">
                {connected ? (
                  <div>
                    <video
                      ref={videoRef}
                      autoPlay
                      playsInline
                      muted
                      controls
                      className="w-full max-w-4xl mx-auto rounded-lg border"
                      style={{display: videoEnabled ? "block" : "none"}}
                    />
                    <div className="mt-2 text-sm text-muted-foreground">
                      <p>
                        Video element ready: {videoRef.current ? "Yes" : "No"}
                      </p>
                      <p>
                        Video srcObject:{" "}
                        {videoRef.current?.srcObject ? "Set" : "Not set"}
                      </p>
                      <p>Video paused: {videoRef.current?.paused}</p>
                      <p>Video readyState: {videoRef.current?.readyState}</p>
                      <p>
                        Video stream available: {videoStream ? "Yes" : "No"}
                      </p>
                      <p>
                        Video tracks:{" "}
                        {videoStream?.getVideoTracks().length || 0}
                      </p>
                    </div>
                  </div>
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
            </Card>
          )}
        </div>
      )}
    </React.Fragment>
  );
};

export default SessionPage;
