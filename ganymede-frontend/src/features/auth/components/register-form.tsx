"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { FormError } from "@/components/form-error";
import { Button } from "@/components/ui/button";
import { registerSchema } from "../schemas/register.schema";
import { useRegister } from "../hooks/use-register";

type FormData = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const router = useRouter();

  const mutation = useRegister();

  const form = useForm<FormData>({
    resolver: zodResolver(registerSchema),
    mode: "onChange",
  });

  function onSubmit(data: FormData) {
    mutation.mutate(data, {
      onSuccess: () => {
        router.push("/login");
      },
    });
  }

  return (
    <main className="min-h-screen flex items-center justify-center">
      <Card className="w-[400px] justify-center ">
        <CardHeader>
          <CardTitle>Create Account</CardTitle>
        </CardHeader>

        <CardContent>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label>Username</Label>

              <Input {...form.register("username")} />

              <FormError message={form.formState.errors.username?.message} />
            </div>

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
              {mutation.isPending ? "Creating..." : "Register"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
