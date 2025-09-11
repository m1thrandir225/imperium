import {Card, CardHeader, CardTitle} from "@/components/ui/card";
import sessionService from "@/services/session.service";
import {useSessionStore} from "@/stores/session.store";
import {Loader2} from "lucide-react";
import React from "react";
import {useEffect, useRef, useState} from "react";
import {useNavigate, useParams, useSearchParams} from "react-router-dom";

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

  const pcRef = useRef<RTCPeerConnection | null>(null);

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
          return;
        }

        if (s.status === "PENDING") {
          setStarting(true);
          const pc = new RTCPeerConnection();
          pcRef.current = pc;
          pc.addTransceiver("video", {direction: "recvonly"});

          const offer = await pc.createOffer();
          await pc.setLocalDescription(offer);

          if (!offer.sdp) {
            throw new Error("Failed to create offer");
          }
          await startSession({webrtc_offer: offer.sdp}, sessionId);

          await sessionService.pollUntilActive(sessionId, (updated) => {
            if (canceled) return;
            setCurrentSession(updated);
          });
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

  return (
    <React.Fragment>
      <Card>
        <CardHeader>
          <CardTitle>Session {sessionId}</CardTitle>
        </CardHeader>
      </Card>

      {(loading || starting) && (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="w-4 h-4 animate-spin" />
        </div>
      )}
      {error && <p className="text-destructive">{error}</p>}

      {!loading && !error && currentSession && (
        <div className="grid gap-3">
          <div className="text-sm text-muted-foreground">
            Status: {currentSession.status}
          </div>
          <div className="text-sm text-muted-foreground">
            Host: {currentSession.host_name} ({hostId})
          </div>
          <div className="text-sm text-muted-foreground">
            Client: {currentSession.client_name} ({clientId})
          </div>
          {/* TODO: host signaling */}
        </div>
      )}
    </React.Fragment>
  );
};

export default SessionPage;
