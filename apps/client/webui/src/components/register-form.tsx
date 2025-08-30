import authService from "@/services/auth.service";
import type {RegisterRequest} from "@/types/responses/auth";
import {zodResolver} from "@hookform/resolvers/zod";
import {useMutation} from "@tanstack/react-query";
import {useForm} from "react-hook-form";
import * as z from "zod";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form";
import {cn} from "@/lib/utils";
import {Input} from "./ui/input";
import {Link, useNavigate} from "react-router-dom";
import {Button} from "./ui/button";
import {Loader2} from "lucide-react";
import {toast} from "sonner";

const formSchema = z.object({
  email: z.email(),
  password: z.string().min(8),
  firstName: z.string().min(1),
  lastName: z.string().min(1),
});

type RegisterFormProps = React.ComponentProps<"div">;

const RegisterForm: React.FC<RegisterFormProps> = (props) => {
  const navigate = useNavigate();
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const {mutateAsync, status} = useMutation({
    mutationKey: ["register"],
    mutationFn: (values: RegisterRequest) => authService.register(values),
    onSuccess: () => {
      toast.success("Account created successfully");
      navigate("/login");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    await mutateAsync({
      email: values.email,
      password: values.password,
      first_name: values.firstName,
      last_name: values.lastName,
    });
  };

  return (
    <div className={cn("flex flex-col gap-6", props.className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome</CardTitle>
          <CardDescription>Create an account to get started</CardDescription>
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
                    name="firstName"
                    render={({field}) => (
                      <FormItem>
                        <FormLabel>First Name</FormLabel>
                        <FormControl>
                          <Input type="text" placeholder="John" {...field} />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                  <FormMessage />
                </div>
                <div className="grid gap-3">
                  <FormField
                    control={form.control}
                    name="lastName"
                    render={({field}) => (
                      <FormItem>
                        <FormLabel>Last Name</FormLabel>
                        <FormControl>
                          <Input type="text" placeholder="Doe" {...field} />
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
                        <FormLabel>Password</FormLabel>
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
                    "Register"
                  )}
                </Button>
              </form>
            </Form>
            <div className="text-center text-sm">
              Already have an account?{" "}
              <Link to="/login" className="underline underline-offset-4">
                Login
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default RegisterForm;
