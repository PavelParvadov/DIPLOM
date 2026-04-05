import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useSession } from "@/app/providers";
import { AdminNav } from "@/features/admin/admin-nav";
import { api } from "@/shared/api/client";
import { formatDate } from "@/shared/lib/format";
import { Button, Card, Field, Input, PageHeader, SecondaryButton } from "@/shared/ui/primitives";

interface InviteCode {
  id: number;
  code: string;
  isActive: boolean;
  createdAt: string;
  expiresAt?: string | null;
  createdByLogin: string;
}

const schema = z.object({
  expiresAt: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

export function AdminInvitesPage() {
  const queryClient = useQueryClient();
  const { selectedHouseId } = useSession();
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { expiresAt: "" },
  });

  const invitesQuery = useQuery({
    queryKey: ["invites", selectedHouseId],
    queryFn: async () => {
      const response = await api.get<{ items: InviteCode[] }>(`/houses/${selectedHouseId}/invite-codes`);
      return response.items;
    },
    enabled: Boolean(selectedHouseId),
  });

  const createMutation = useMutation({
    mutationFn: (values: FormValues) =>
      api.post(`/houses/${selectedHouseId}/invite-codes`, {
        expiresAt: values.expiresAt ? new Date(values.expiresAt).toISOString() : null,
      }),
    onSuccess() {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["invites", selectedHouseId] });
    },
  });

  const deactivateMutation = useMutation({
    mutationFn: (id: number) => api.patch(`/houses/${selectedHouseId}/invite-codes/${id}/deactivate`),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["invites", selectedHouseId] });
    },
  });

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Admin"
        title="Invite codes"
        description="Создавайте коды приглашения для новых соседей и отключайте их после использования."
        actions={<AdminNav />}
      />
      <Card className="space-y-5">
        <form className="flex flex-col gap-4 md:flex-row md:items-end" noValidate onSubmit={form.handleSubmit((values) => createMutation.mutate(values))}>
          <div className="flex-1">
            <Field label="Дата окончания действия" error={form.formState.errors.expiresAt?.message}>
              <Input {...form.register("expiresAt")} type="datetime-local" />
            </Field>
          </div>
          <Button type="submit" disabled={createMutation.isPending}>
            Сгенерировать код
          </Button>
        </form>
      </Card>
      <div className="grid gap-4">
        {invitesQuery.data?.map((invite) => (
          <Card key={invite.id} className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="space-y-1">
              <div className="font-mono text-lg font-semibold text-ink">{invite.code}</div>
              <div className="text-sm text-slate-500">
                Создан: {formatDate(invite.createdAt)} · Автор: {invite.createdByLogin}
              </div>
              <div className="text-sm text-slate-500">
                Статус: {invite.isActive ? "активен" : "выключен"}
                {invite.expiresAt ? ` · Действует до ${formatDate(invite.expiresAt)}` : ""}
              </div>
            </div>
            {invite.isActive ? <SecondaryButton onClick={() => deactivateMutation.mutate(invite.id)}>Деактивировать</SecondaryButton> : null}
          </Card>
        ))}
      </div>
    </div>
  );
}
