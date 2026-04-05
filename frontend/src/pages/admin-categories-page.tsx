import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useSession } from "@/app/providers";
import type { Category } from "@/entities/category/model";
import { AdminNav } from "@/features/admin/admin-nav";
import { api } from "@/shared/api/client";
import { Button, Card, Field, Input, PageHeader, SecondaryButton } from "@/shared/ui/primitives";

const schema = z.object({
  name: z.string().min(2, "Минимум 2 символа"),
});

type FormValues = z.infer<typeof schema>;

export function AdminCategoriesPage() {
  const queryClient = useQueryClient();
  const { selectedHouseId } = useSession();
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: "" },
  });

  const categoriesQuery = useQuery({
    queryKey: ["categories", selectedHouseId],
    queryFn: async () => {
      const response = await api.get<{ items: Category[] }>(`/houses/${selectedHouseId}/categories`);
      return response.items;
    },
    enabled: Boolean(selectedHouseId),
  });

  const createMutation = useMutation({
    mutationFn: (values: FormValues) => api.post(`/houses/${selectedHouseId}/categories`, values),
    onSuccess() {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["categories", selectedHouseId] });
    },
  });

  const renameMutation = useMutation({
    mutationFn: (payload: { id: number; name: string }) =>
      api.patch(`/houses/${selectedHouseId}/categories/${payload.id}`, { name: payload.name }),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["categories", selectedHouseId] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/houses/${selectedHouseId}/categories/${id}`),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["categories", selectedHouseId] });
    },
  });

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Admin"
        title="Управление категориями"
        description="Категории задают структуру ленты и помогают держать контент дома аккуратным и понятным."
        actions={<AdminNav />}
      />

      <Card className="space-y-5">
        <form className="flex flex-col gap-4 md:flex-row md:items-end" noValidate onSubmit={form.handleSubmit((values) => createMutation.mutate(values))}>
          <div className="flex-1">
            <Field label="Новая категория" error={form.formState.errors.name?.message}>
              <Input {...form.register("name")} placeholder="Благоустройство" />
            </Field>
          </div>
          <Button type="submit" disabled={createMutation.isPending}>
            Добавить
          </Button>
        </form>
      </Card>

      <div className="grid gap-4">
        {categoriesQuery.data?.map((category) => (
          <Card key={category.id} className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <div className="text-lg font-semibold text-ink">{category.name}</div>
            </div>
            <div className="flex flex-wrap gap-3">
              <SecondaryButton
                onClick={() => {
                  const nextName = window.prompt("Новое название категории", category.name);
                  if (nextName && nextName !== category.name) {
                    renameMutation.mutate({ id: category.id, name: nextName });
                  }
                }}
              >
                Переименовать
              </SecondaryButton>
              <SecondaryButton onClick={() => deleteMutation.mutate(category.id)}>Удалить</SecondaryButton>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
