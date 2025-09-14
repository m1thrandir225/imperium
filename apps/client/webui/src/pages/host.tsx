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

  useEffect(() => {
    if (!hostId || hostId === "") {
      navigate("/hosts", {replace: true});
    }
  }, [hostId, navigate]);

  const {
    data: host,
    error,
    isLoading,
  } = useQuery({
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
    enabled: !!hostId && !!host,
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

  return (
    <React.Fragment>
      {isLoading && (
        <div className="flex justify-center items-center h-full">
          <Loader2 className="w-4 h-4 animate-spin" />
        </div>
      )}
      {host && !isLoading && (
        <React.Fragment>
          <Card>
            <CardHeader>
              <CardTitle>Host {host?.name}</CardTitle>
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
            {programs &&
              !programsLoading &&
              programs.map((program) => (
                <ProgramItem key={program.id} program={program} host={host} />
              ))}
          </div>
        </React.Fragment>
      )}
    </React.Fragment>
  );
};

export default SingleHostPage;
