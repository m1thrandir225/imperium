import ProgramItem from "@/components/program-item";
import {Card, CardHeader, CardTitle} from "@/components/ui/card";
import hostService from "@/services/host.service";
import {useQuery} from "@tanstack/react-query";
import {isAxiosError} from "axios";
import {Loader2} from "lucide-react";
import React, {useEffect} from "react";
import {useNavigate, useParams} from "react-router-dom";

const SingleHostPage: React.FC = () => {
  const {hostId} = useParams();
  const navigate = useNavigate();

  // const videoRef = useRef<HTMLVideoElement>(null);
  // const [connecting, setConnecting] = useState(false);
  // const [pc, setPc] = useState<RTCPeerConnection | null>(null);
  // const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (!hostId || hostId === "") {
      navigate("/hosts", {replace: true});
    }
  }, [hostId, navigate]);

  const {data, error, isLoading} = useQuery({
    queryKey: ["host", hostId],
    queryFn: () => hostService.getHost(hostId as string),
    enabled: !!hostId,
  });

  const {
    data: programs,
    error: programsError,
    isLoading: programsLoading,
  } = useQuery({
    queryKey: ["programs", hostId],
    queryFn: () => hostService.getHostPrograms(hostId as string),
    enabled: !!hostId && !!data,
  });

  useEffect(() => {
    if (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        navigate("/hosts", {replace: true});
      } else if (error.message && error.message.includes("404")) {
        navigate("/hosts", {replace: true});
      }
    }
  }, [error, navigate]);

  // const connect = async () => {
  //   if (!data) return;
  //   setConnecting(true);
  //   try {
  //     const hostIP = data.ip_address;

  //     // TODO: actuall sessionID from auth-server needs to be used
  //     const sessionId = String(hostId);

  //     const signalingURL = `http://${hostIP}:8090/api/session/webrtc/offer`;
  //     const wsURL = `ws://${hostIP}:8080/ws?session_id=${sessionId}`;

  //     const peer = new RTCPeerConnection();
  //     peer.addTransceiver("video", {direction: "recvonly"});
  //     peer.ontrack = (ev) => {
  //       if (videoRef.current) videoRef.current.srcObject = ev.streams[0];
  //     };

  //     const offer = await peer.createOffer();
  //     await peer.setLocalDescription(offer);

  //     const res = await fetch(signalingURL, {
  //       method: "POST",
  //       headers: {"Content-Type": "application/json"},
  //       body: JSON.stringify({sdp: offer.sdp}),
  //     });
  //     if (!res.ok) throw new Error("Signaling failed");
  //     const answer = await res.json();
  //     await peer.setRemoteDescription({type: "answer", sdp: answer.sdp});

  //     setPc(peer);

  //     const socket = new WebSocket(wsURL);
  //     socket.onopen = () => console.log("WS connected");
  //     socket.onclose = () => console.log("WS closed");
  //     setWs(socket);
  //   } catch (e) {
  //     console.error(e);
  //   } finally {
  //     setConnecting(false);
  //   }
  // };

  // useEffect(() => {
  //   const v = videoRef.current;
  //   if (!v) return;

  //   const onMouseDown = (e: MouseEvent) => {
  //     const payload = {
  //       type: "mouse",
  //       action: "click",
  //       button: e.button === 0 ? "left" : "right",
  //     };
  //     ws?.send(JSON.stringify(payload));
  //   };
  //   const onKeyDown = (e: KeyboardEvent) =>
  //     ws?.send(JSON.stringify({type: "keyboard", action: "press", key: e.key}));
  //   const onKeyUp = (e: KeyboardEvent) =>
  //     ws?.send(
  //       JSON.stringify({type: "keyboard", action: "release", key: e.key})
  //     );

  //   v.addEventListener("mousedown", onMouseDown);
  //   window.addEventListener("keydown", onKeyDown);
  //   window.addEventListener("keyup", onKeyUp);

  //   return () => {
  //     v.removeEventListener("mousedown", onMouseDown);
  //     window.removeEventListener("keydown", onKeyDown);
  //     window.removeEventListener("keyup", onKeyUp);
  //     pc?.close();
  //     ws?.close();
  //   };
  // }, [pc, ws]);

  return (
    <React.Fragment>
      {isLoading && (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="w-4 h-4 animate-spin" />
        </div>
      )}
      {data && !isLoading && (
        <React.Fragment>
          <Card>
            <CardHeader>
              <CardTitle>Host {data?.name}</CardTitle>
            </CardHeader>
          </Card>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 w-full h-full">
            {programsLoading && (
              <div className="flex justify-center items-center h-full col-span-full">
                <Loader2 className="w-4 h-4 animate-spin" />
              </div>
            )}
            {programsError && (
              <p className="col-span-full">Error: {programsError.message}</p>
            )}
            {programs && !programsLoading && (
              <div className="grid grid-cols-2 gap-3">
                {programs.map((program) => (
                  <ProgramItem key={program.id} program={program} />
                ))}
              </div>
            )}
          </div>
        </React.Fragment>
      )}
    </React.Fragment>
  );
};

export default SingleHostPage;
