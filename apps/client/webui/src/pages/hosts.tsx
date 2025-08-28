import HostItem from "@/components/host-item";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import hostService from "@/services/host.service";
import {useQuery} from "@tanstack/react-query";
import {Loader2} from "lucide-react";
import React from "react";

const HostsPage: React.FC = () => {
  const {data, isLoading, error} = useQuery({
    queryKey: ["hosts"],
    queryFn: () => hostService.getHosts(),
  });
  return (
    <React.Fragment>
      <Card>
        <CardHeader>
          <CardTitle>Hosts</CardTitle>
          <CardDescription>
            This is a list of hosts that you can connect to.
          </CardDescription>
        </CardHeader>
      </Card>
      {isLoading && (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="w-4 h-4 animate-spin" />
        </div>
      )}
      {error && <p>Error: {error.message}</p>}
      <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-3">
        {data && data.map((host) => <HostItem key={host.id} {...host} />)}
      </div>
    </React.Fragment>
  );
};

export default HostsPage;
