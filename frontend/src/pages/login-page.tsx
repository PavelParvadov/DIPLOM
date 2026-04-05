import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";

import { useSession } from "@/app/providers";
import type { User } from "@/entities/user/model";
import { api } from "@/shared/api/client";
import { Button, Card, Field, Input } from "@/shared/ui/primitives";

const schema = z.object({
  login: z.string().min(3, "Минимум 3 символа"),
  password: z.string().min(6, "Минимум 6 символов"),
});

type FormValues = z.infer<typeof schema>;
type AuthPayload = { user: User; tokens: { accessToken: string; refreshToken: string; expiresAt: string } };

export function LoginPage() {
  const navigate = useNavigate();
  const { setAuth } = useSession();
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { login: "", password: "" },
  });

  const mutation = useMutation({
    mutationFn: (values: FormValues) => api.post<AuthPayload>("/auth/login", values),
    onSuccess(payload: AuthPayload) {
      setAuth(payload.user, payload.tokens);
      navigate("/posts");
    },
  });
  const errorMessage = mutation.error instanceof Error ? mutation.error.message : "";

  return (
    <div className="flex min-h-screen items-center justify-center px-4 py-10">
      <Card className="grid max-w-5xl overflow-hidden p-0 md:grid-cols-[1.1fr_0.9fr]">
        <div className="bg-ink px-8 py-10 text-white md:px-10">
          <div className="space-y-5">
            <div className="text-xs uppercase tracking-[0.3em] text-sand/80">HappyHouse</div>
            <h1 className="font-display text-4xl">Общение соседей без хаоса в мессенджерах.</h1>
            <p className="max-w-md text-sm leading-6 text-slate-300">
              Публикуйте объявления, новости дома, обсуждайте важные темы и управляйте категориями в одном аккуратном веб-интерфейсе.
            </p>
          </div>
        </div>
        <div className="px-8 py-10 md:px-10">
          <div className="mb-8 space-y-2">
            <h2 className="text-3xl font-semibold text-ink">Вход</h2>
            <p className="text-sm text-slate-500">Введите логин и пароль, чтобы открыть ленту дома.</p>
          </div>
          <form className="space-y-5" noValidate onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
            <Field label="Логин" error={form.formState.errors.login?.message}>
              <Input {...form.register("login")} placeholder="admin_demo" />
            </Field>
            <Field label="Пароль" error={form.formState.errors.password?.message}>
              <Input {...form.register("password")} type="password" placeholder="••••••••" />
            </Field>
            {errorMessage ? <div className="text-sm text-red-600">{errorMessage}</div> : null}
            <Button className="w-full" type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Входим..." : "Войти"}
            </Button>
          </form>
          <p className="mt-6 text-sm text-slate-500">
            Нет аккаунта?{" "}
            <Link to="/register" className="font-semibold text-moss">
              Зарегистрироваться
            </Link>
          </p>
        </div>
      </Card>
    </div>
  );
}
