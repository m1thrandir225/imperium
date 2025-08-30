import Logo from "@/components/logo";
import {Button} from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {Toaster} from "@/components/ui/sonner";
import configService from "@/services/config.service";
import useConfigStore from "@/stores/config.store";
import type {SetupConfigRequest} from "@/types/responses/config";
import {zodResolver} from "@hookform/resolvers/zod";
import {useMutation} from "@tanstack/react-query";
import {Loader2} from "lucide-react";
import type React from "react";
import {useForm} from "react-hook-form";
import {toast} from "sonner";
import * as z from "zod";

const formSchema = z.object({
  authServerBaseUrl: z.url(),
});

const SetupPage: React.FC = () => {
  const setConfigured = useConfigStore((state) => state.setupConfiguration);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const {mutateAsync, status} = useMutation({
    mutationKey: ["setup"],
    mutationFn: (values: SetupConfigRequest) =>
      configService.setupConfig(values),
    onSuccess: () => {
      toast.success("Config setup successfully");
      setConfigured(true);
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    await mutateAsync({
      auth_server_base_url: values.authServerBaseUrl.toString(),
    });
  };

  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Logo variant="default" />
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Welcome to Imperium</CardTitle>
            <CardDescription>
              Let's set up your application configuration to get started.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                className="grid gap-6"
                onSubmit={form.handleSubmit(onSubmit)}
              >
                <FormField
                  control={form.control}
                  name="authServerBaseUrl"
                  render={({field}) => (
                    <FormItem>
                      <FormLabel> Authentication Server Base URL</FormLabel>
                      <FormControl>
                        <Input
                          type="url"
                          placeholder="https://auth.example.com"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" disabled={status === "pending"}>
                  {status === "pending" ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    "Continue"
                  )}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>

      <Toaster />
    </div>
  );
};

export default SetupPage;
