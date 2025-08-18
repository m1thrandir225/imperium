import {LoginForm} from "@/components/login-form";
import AuthLayout from "@/layouts/auth-layout";
import type React from "react";

const LoginPage: React.FC = () => {
  return (
    <AuthLayout>
      <LoginForm />
    </AuthLayout>
  );
};

export default LoginPage;
