import type {Program} from "@/types/models/program";
import type React from "react";
import {Card, CardFooter, CardHeader, CardTitle} from "./ui/card";
import {Button} from "./ui/button";
import {useNavigate} from "react-router-dom";
import useAuthStore from "@/stores/auth.store";
import {useState} from "react";
import sessionService from "@/services/session.service";
import {Loader2} from "lucide-react";
import type {Host} from "@/types/models/host";
import {useSessionStore} from "@/stores/session.store";

type ComponentProps = {
  program: Program;
  host: Host;
};

const ProgramItem: React.FC<ComponentProps> = (props) => {
  const {program, host} = props;
  const navigate = useNavigate();
  const client = useAuthStore((state) => state.client);
  const sessionStore = useSessionStore();

  const [isCreating, setIsCreating] = useState(false);

  const handleStartSession = async () => {
    if (!client?.id) {
      console.error("Client not found");
      return;
    }
    setIsCreating(true);
    try {
      const session = await sessionService.create({
        host_id: host.id,
        client_id: client.id,
        program_id: props.program.id,
      });
      sessionStore.setHost(host);
      sessionStore.setCurrentSession(session);

      navigate(
        `/sessions/${session.id}?host_id=${host.id}&client_id=${client.id}`
      );
    } catch (e: unknown) {
      console.error(e);
    } finally {
      setIsCreating(false);
    }
  };
  return (
    <Card>
      <CardHeader>
        <CardTitle>{program.name}</CardTitle>
      </CardHeader>
      <CardFooter>
        <Button
          onClick={handleStartSession}
          disabled={isCreating || !client?.id}
          className="w-full"
        >
          {isCreating ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            "Start Session"
          )}
        </Button>
      </CardFooter>
    </Card>
  );
};

export default ProgramItem;
