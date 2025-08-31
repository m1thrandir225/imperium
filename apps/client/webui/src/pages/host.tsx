import {Card, CardHeader, CardTitle} from "@/components/ui/card";
import hostService from "@/services/host.service";
import {useQuery} from "@tanstack/react-query";
import {isAxiosError} from "axios";
import {Loader2} from "lucide-react";
import React from "react";
import {useEffect} from "react";
import {useNavigate, useParams} from "react-router-dom";

const SingleHostPage: React.FC = () => {
  const {hostId} = useParams();
  const navigate = useNavigate();

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

  useEffect(() => {
    if (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
        navigate("/hosts", {replace: true});
      }
      // Fallback: check if error message contains 404 (for backwards compatibility)
      else if (error.message && error.message.includes("404")) {
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
      {data && !isLoading && (
        <>
          <Card>
            <CardHeader>
              <CardTitle>Host {data?.name}</CardTitle>
            </CardHeader>
          </Card>
        </>
      )}
    </React.Fragment>
  );
};

export default SingleHostPage;
