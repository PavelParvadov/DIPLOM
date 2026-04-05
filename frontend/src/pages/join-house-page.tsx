import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { z } from "zod";

import { useSession } from "@/app/providers";
import { api } from "@/shared/api/client";
import { Button, Card, Field, Input, PageHeader } from "@/shared/ui/primitives";

const joinSchema = z.object({
  code: z.string().min(4, "Введите код приглашения"),
});

const createHouseSchema = z.object({
  name: z.string().min(3, "Минимум 3 символа"),
  address: z.string().min(5, "Минимум 5 символов"),
});

type JoinFormValues = z.infer<typeof joinSchema>;
type CreateHouseFormValues = z.infer<typeof createHouseSchema>;

export function JoinHousePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { setSelectedHouseId } = useSession();
  const joinForm = useForm<JoinFormValues>({
    resolver: zodResolver(joinSchema),
    defaultValues: { code: "" },
  });
  const createHouseForm = useForm<CreateHouseFormValues>({
    resolver: zodResolver(createHouseSchema),
    defaultValues: { name: "", address: "" },
  });

  const joinMutation = useMutation({
    mutationFn: (values: JoinFormValues) => api.post<{ success: boolean }>("/houses/join-by-code", values),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["houses"] });
      navigate("/posts");
    },
  });

  const createHouseMutation = useMutation({
    mutationFn: (values: CreateHouseFormValues) => api.post<{ item: { id: number } }>("/houses", values),
    onSuccess(response) {
      createHouseForm.reset();
      setSelectedHouseId(response.item.id);
      queryClient.invalidateQueries({ queryKey: ["houses"] });
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      navigate("/posts");
    },
  });

  const joinErrorMessage = joinMutation.error instanceof Error ? joinMutation.error.message : "";
  const createHouseErrorMessage = createHouseMutation.error instanceof Error ? createHouseMutation.error.message : "";

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Дома"
        title="Создать дом или присоединиться по коду"
        description="Любой пользователь может открыть свой дом и автоматически стать его администратором. Если дом уже существует, используйте invite code от администратора."
      />

      <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
        <Card className="space-y-6">
          <div className="space-y-2">
            <h2 className="text-2xl font-semibold text-ink">Создать новый дом</h2>
            <p className="text-sm leading-6 text-slate-500">После создания вы сразу получите роль администратора и стандартные категории для первых публикаций.</p>
          </div>
          <form className="space-y-5" noValidate onSubmit={createHouseForm.handleSubmit((values) => createHouseMutation.mutate(values))}>
            <Field label="Название дома" error={createHouseForm.formState.errors.name?.message}>
              <Input {...createHouseForm.register("name")} placeholder="ЖК Сосновый двор" />
            </Field>
            <Field label="Адрес" error={createHouseForm.formState.errors.address?.message}>
              <Input {...createHouseForm.register("address")} placeholder="Санкт-Петербург, Невский проспект, 18" />
            </Field>
            {createHouseErrorMessage ? <div className="text-sm text-red-600">{createHouseErrorMessage}</div> : null}
            <Button type="submit" disabled={createHouseMutation.isPending}>
              {createHouseMutation.isPending ? "Создаем дом..." : "Создать дом"}
            </Button>
          </form>
        </Card>

        <Card className="space-y-6">
          <div className="space-y-2">
            <h2 className="text-2xl font-semibold text-ink">Вступить по invite code</h2>
            <p className="text-sm leading-6 text-slate-500">Если администратор уже создал дом, достаточно ввести код приглашения. При повторном вступлении покажем понятное сообщение вместо технической ошибки.</p>
          </div>
          <form className="space-y-5" noValidate onSubmit={joinForm.handleSubmit((values) => joinMutation.mutate(values))}>
            <Field label="Invite code" error={joinForm.formState.errors.code?.message}>
              <Input {...joinForm.register("code")} placeholder="HAPPY2026" />
            </Field>
            {joinErrorMessage ? <div className="text-sm text-red-600">{joinErrorMessage}</div> : null}
            <Button type="submit" disabled={joinMutation.isPending}>
              {joinMutation.isPending ? "Подключаем..." : "Вступить в дом"}
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
}
