import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useId, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { z } from "zod";

import { useSession } from "@/app/providers";
import type { Category } from "@/entities/category/model";
import { api } from "@/shared/api/client";
import { Button, Card, Field, Input, PageHeader, Select, Textarea } from "@/shared/ui/primitives";

const schema = z.object({
  categoryId: z.coerce.number().min(1, "Выберите категорию"),
  title: z.string().min(3, "Минимум 3 символа"),
  content: z.string().min(5, "Минимум 5 символов"),
  image: z.any().optional(),
});

type FormValues = z.infer<typeof schema>;

export function NewPostPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { selectedHouseId } = useSession();
  const inputId = useId();
  const [previewUrl, setPreviewUrl] = useState("");
  const categoriesQuery = useQuery({
    queryKey: ["categories", selectedHouseId],
    queryFn: async () => {
      const response = await api.get<{ items: Category[] }>(`/houses/${selectedHouseId}/categories`);
      return response.items;
    },
    enabled: Boolean(selectedHouseId),
  });

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { categoryId: 0, title: "", content: "" },
  });

  const watchedImage = form.watch("image") as FileList | undefined;
  const selectedFile = watchedImage?.[0];
  useEffect(() => {
    if (!selectedFile) {
      setPreviewUrl("");
      return;
    }
    const objectUrl = URL.createObjectURL(selectedFile);
    setPreviewUrl(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [selectedFile]);

  const mutation = useMutation({
    mutationFn: (values: FormValues) => {
      const payload = new FormData();
      payload.append("categoryId", String(values.categoryId));
      payload.append("title", values.title);
      payload.append("content", values.content);
      const imageFiles = values.image as FileList | undefined;
      if (imageFiles?.[0]) {
        payload.append("image", imageFiles[0]);
      }
      return api.post(`/houses/${selectedHouseId}/posts`, payload);
    },
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      navigate("/posts");
    },
  });
  const errorMessage = mutation.error instanceof Error ? mutation.error.message : "";

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Новый контент"
        title="Создать пост"
        description="Напишите короткий и понятный пост, добавьте при необходимости изображение и отправьте его в ленту дома."
      />
      <Card className="max-w-3xl">
        <form className="space-y-5" noValidate onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
          <Field label="Категория" error={form.formState.errors.categoryId?.message}>
            <Select {...form.register("categoryId")}>
              <option value={0}>Выберите категорию</option>
              {categoriesQuery.data?.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </Select>
          </Field>
          <Field label="Заголовок" error={form.formState.errors.title?.message}>
            <Input {...form.register("title")} placeholder="Собираем предложения по благоустройству двора" />
          </Field>
          <Field label="Содержание" error={form.formState.errors.content?.message}>
            <Textarea {...form.register("content")} placeholder="Опишите идею, сроки, детали и что требуется от соседей." />
          </Field>

          <div className="flex items-center gap-3">
            <input {...form.register("image")} id={inputId} type="file" accept="image/png,image/jpeg,image/webp,image/gif" className="hidden" />
            <label
              htmlFor={inputId}
              className="inline-flex h-12 w-12 cursor-pointer items-center justify-center rounded-2xl border border-slate-200 bg-white text-ink transition hover:border-moss hover:text-moss"
              title="Прикрепить изображение"
            >
              <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.9" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21.44 11.05 12 20.5a6 6 0 0 1-8.49-8.49l10.6-10.61a4 4 0 0 1 5.66 5.66L9.17 17.66a2 2 0 0 1-2.83-2.83l9.2-9.19" />
              </svg>
            </label>
            <div className="text-sm text-slate-500">
              {selectedFile ? `Выбрано: ${selectedFile.name}` : "Прикрепить изображение"}
            </div>
          </div>

          {previewUrl ? (
            <div className="overflow-hidden rounded-[24px] border border-slate-200 bg-slate-50">
              <img src={previewUrl} alt="Предпросмотр изображения" className="h-72 w-full object-cover" />
            </div>
          ) : null}
          {errorMessage ? <div className="text-sm text-red-600">{errorMessage}</div> : null}
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Публикуем..." : "Опубликовать пост"}
          </Button>
        </form>
      </Card>
    </div>
  );
}
