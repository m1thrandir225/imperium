import type {SimpleHost} from "@/types/models/host";
import type React from "react";
import {
  Card,
  CardAction,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {Button} from "./ui/button";
import {useNavigate} from "react-router-dom";

type ComponentProps = {
  id: string;
  name: string;
  status: SimpleHost["status"];
};

const HostItem: React.FC<ComponentProps> = (props) => {
  const navigate = useNavigate();

  const handleConnect = () => {
    navigate(`/hosts/${props.id}`);
  };
  return (
    <Card className="w-full h-full">
      <CardHeader>
        <CardTitle>{props.name}</CardTitle>
        <CardDescription>{props.status}</CardDescription>
        <CardAction>
          <Button
            variant={"outline"}
            disabled={!props.status.includes("AVAILABLE")}
            onClick={handleConnect}
          >
            {props.status === "AVAILABLE" ? "Connect" : "Not Available"}
          </Button>
        </CardAction>
      </CardHeader>
    </Card>
  );
};

export default HostItem;
