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
  displayName: z.string().min(2, "Введите имя для отображения"),
  login: z.string().min(3, "Минимум 3 символа"),
  password: z.string().min(6, "Минимум 6 символов"),
});

type FormValues = z.infer<typeof schema>;
type AuthPayload = { user: User; tokens: { accessToken: string; refreshToken: string; expiresAt: string } };

export function RegisterPage() {
  const navigate = useNavigate();
  const { setAuth } = useSession();
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { displayName: "", login: "", password: "" },
  });

  const mutation = useMutation({
    mutationFn: (values: FormValues) => api.post<AuthPayload>("/auth/register", values),
    onSuccess(payload: AuthPayload) {
      setAuth(payload.user, payload.tokens);
      navigate("/join");
    },
  });
  const errorMessage = mutation.error instanceof Error ? mutation.error.message : "";

  return (
    <div className="flex min-h-screen items-center justify-center px-4 py-10">
      <Card className="w-full max-w-xl space-y-8">
        <div className="space-y-2">
          <div className="text-xs uppercase tracking-[0.3em] text-moss">Новый аккаунт</div>
          <h1 className="font-display text-4xl text-ink">Создать профиль в HappyHouse</h1>
          <p className="text-sm text-slate-500">После регистрации вы сможете вступить в дом по invite code и сразу пользоваться лентой.</p>
        </div>
        <form className="space-y-5" noValidate onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
          <Field label="Имя для отображения" error={form.formState.errors.displayName?.message}>
            <Input {...form.register("displayName")} placeholder="Павел Петров" />
          </Field>
          <Field label="Логин" error={form.formState.errors.login?.message}>
            <Input {...form.register("login")} placeholder="pavel_petrov" />
          </Field>
          <Field label="Пароль" error={form.formState.errors.password?.message}>
            <Input {...form.register("password")} type="password" placeholder="••••••••" />
          </Field>
          {errorMessage ? <div className="text-sm text-red-600">{errorMessage}</div> : null}
          <Button className="w-full" type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Создаем аккаунт..." : "Создать аккаунт"}
          </Button>
        </form>
        <p className="text-sm text-slate-500">
          Уже есть аккаунт?{" "}
          <Link to="/login" className="font-semibold text-moss">
            Перейти ко входу
          </Link>
        </p>
      </Card>
    </div>
  );
}
