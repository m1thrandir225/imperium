import type {Program} from "@/types/models/program";
import type React from "react";
import {Card, CardHeader, CardTitle} from "./ui/card";

type ComponentProps = {
  program: Program;
};

const ProgramItem: React.FC<ComponentProps> = (props) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{props.program.name}</CardTitle>
      </CardHeader>
    </Card>
  );
};

export default ProgramItem;
