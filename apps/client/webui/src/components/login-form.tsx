import {Button} from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {Input} from "@/components/ui/input";
import {cn} from "@/lib/utils";
import {zodResolver} from "@hookform/resolvers/zod";
import {useForm} from "react-hook-form";
import {Link, useNavigate} from "react-router-dom";
import * as z from "zod";
import {
  FormControl,
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form";
import {useMutation} from "@tanstack/react-query";
import type {LoginRequest} from "@/types/responses/auth";
import authService from "@/services/auth.service";
import {Loader2} from "lucide-react";
import useAuthStore from "@/stores/auth.store";

const formSchema = z.object({
  email: z.email(),
  password: z.string().min(8),
});

export function LoginForm({className, ...props}: React.ComponentProps<"div">) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const authStore = useAuthStore();
  const navigate = useNavigate();

  const {mutateAsync, status} = useMutation({
    mutationKey: ["login"],
    mutationFn: (values: LoginRequest) => authService.login(values),
    onSuccess: async (data) => {
      await authStore.login(data);
      navigate("/hosts");
    },
    onError: (error) => {
      console.error(error.message);
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    await mutateAsync({
      email: values.email,
      password: values.password,
    });
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>Login with your email and password</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6">
            <Form {...form}>
              <form
                className="grid gap-6"
                onSubmit={form.handleSubmit(onSubmit)}
              >
                <div className="grid gap-3">
                  <FormField
                    control={form.control}
                    name="email"
                    render={({field}) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <Input
                            type="email"
                            placeholder="you@example.com"
                            {...field}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                  <FormMessage />
                </div>
                <div className="grid gap-3">
                  <FormField
                    control={form.control}
                    name="password"
                    render={({field}) => (
                      <FormItem>
                        <div className="flex items-center">
                          <FormLabel>Password</FormLabel>
                          <Link
                            to={"/forgot-password"}
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                          >
                            Forgot your password?
                          </Link>
                        </div>
                        <FormControl>
                          <Input
                            type="password"
                            placeholder="**********"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <Button
                  disabled={status === "pending"}
                  type="submit"
                  className="w-full"
                >
                  {status === "pending" ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    "Login"
                  )}
                </Button>
              </form>
            </Form>
            <div className="text-center text-sm">
              Don&apos;t have an account?{" "}
              <Link to="/register" className="underline underline-offset-4">
                Sign up
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
