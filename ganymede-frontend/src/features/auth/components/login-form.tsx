"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { z } from "zod";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { FormError } from "@/components/form-error";
import { Button } from "@/components/ui/button";
import { useLogin } from "../hooks/use-login";
import { loginSchema } from "../schemas/login.schema";

type FormData = z.infer<typeof loginSchema>;

export function LoginForm() {
  const router = useRouter();
  const mutation = useLogin();

  const form = useForm<FormData>({
    resolver: zodResolver(loginSchema),
    mode: "onChange",
  });

  function onSubmit(data: FormData) {
    mutation.mutate(data, {
      onSuccess: () => {
        router.push("/dashboard");
      },
    });
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-200">
      <Card className="w-[500px]">
        <CardHeader>
          <CardTitle>Login</CardTitle>
        </CardHeader>

        <CardContent>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label>Email</Label>
              <Input {...form.register("email")} />
              <FormError message={form.formState.errors.email?.message} />
            </div>

            <div>
              <Label>Password</Label>
              <Input type="password" {...form.register("password")} />
              <FormError message={form.formState.errors.password?.message} />
            </div>

            <Button
              className="w-full"
              type="submit"
              disabled={mutation.isPending}
            >
              {mutation.isPending ? "Logging in..." : "Login"}
            </Button>

            <Button className="w-full" onClick={() => router.push("/register")}>
              Register
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
