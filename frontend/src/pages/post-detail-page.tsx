import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router-dom";
import { z } from "zod";

import { useSession } from "@/app/providers";
import type { CommentListResponse } from "@/entities/comment/model";
import type { Post } from "@/entities/post/model";
import { CommentList } from "@/features/comments/comment-list";
import { api, toAssetUrl } from "@/shared/api/client";
import { formatDate } from "@/shared/lib/format";
import { Button, Card, Field, PageHeader, SecondaryButton, Textarea } from "@/shared/ui/primitives";

const schema = z.object({
  content: z.string().min(1, "Введите комментарий"),
});

type FormValues = z.infer<typeof schema>;

export function PostDetailPage() {
  const { postId } = useParams();
  const queryClient = useQueryClient();
  const { selectedHouseId } = useSession();
  const [commentsPage, setCommentsPage] = useState(1);
  const pageSize = 4;
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { content: "" },
  });

  const postQuery = useQuery({
    queryKey: ["post", selectedHouseId, postId],
    queryFn: async () => {
      const response = await api.get<{ item: Post }>(`/houses/${selectedHouseId}/posts/${postId}`);
      return response.item;
    },
    enabled: Boolean(selectedHouseId && postId),
  });

  const commentsQuery = useQuery({
    queryKey: ["comments", selectedHouseId, postId, commentsPage],
    queryFn: () =>
      api.get<CommentListResponse>(`/houses/${selectedHouseId}/posts/${postId}/comments?page=${commentsPage}&pageSize=${pageSize}`),
    enabled: Boolean(selectedHouseId && postId),
  });

  const mutation = useMutation({
    mutationFn: (values: FormValues) => api.post(`/houses/${selectedHouseId}/posts/${postId}/comments`, values),
    onSuccess() {
      form.reset();
      setCommentsPage(1);
      queryClient.invalidateQueries({ queryKey: ["comments", selectedHouseId, postId] });
    },
  });

  const errorMessage = mutation.error instanceof Error ? mutation.error.message : "";
  const post = postQuery.data;
  const comments = commentsQuery.data?.items ?? [];
  const totalComments = commentsQuery.data?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(totalComments / pageSize));

  return (
    <div className="space-y-6">
      {post ? (
        <PageHeader
          eyebrow={post.categoryName}
          title={post.title}
          description={`Автор: ${post.authorName} · ${formatDate(post.createdAt)}`}
        />
      ) : null}

      {post ? (
        <Card className="space-y-5">
          {post.imageUrl ? (
            <div className="overflow-hidden rounded-[24px] border border-slate-200 bg-slate-50">
              <img src={toAssetUrl(post.imageUrl)} alt={post.title} className="max-h-[32rem] w-full object-cover" />
            </div>
          ) : null}
          <p className="text-base leading-7 text-slate-700">{post.content}</p>
        </Card>
      ) : null}

      <div className="grid gap-6 xl:grid-cols-[1fr_0.9fr]">
        <Card className="space-y-5">
          <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div className="min-w-0">
              <h2 className="text-xl font-semibold text-ink">Комментарии</h2>
              <p className="text-sm text-slate-500">Обсуждение вынесено в отдельный компактный блок, чтобы не перегружать экран.</p>
            </div>
            <span className="self-start whitespace-nowrap rounded-full bg-sand px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-ink">
              {totalComments} всего
            </span>
          </div>

          {comments.length ? (
            <>
              <CommentList comments={comments} />
              {totalPages > 1 ? (
                <div className="flex flex-col gap-3 border-t border-slate-100 pt-4 sm:flex-row sm:items-center sm:justify-between">
                  <span className="text-sm text-slate-500">
                    Страница {commentsPage} из {totalPages}
                  </span>
                  <div className="flex gap-3">
                    <SecondaryButton disabled={commentsPage === 1} onClick={() => setCommentsPage((page) => Math.max(1, page - 1))}>
                      Назад
                    </SecondaryButton>
                    <SecondaryButton
                      disabled={commentsPage >= totalPages}
                      onClick={() => setCommentsPage((page) => Math.min(totalPages, page + 1))}
                    >
                      Еще
                    </SecondaryButton>
                  </div>
                </div>
              ) : null}
            </>
          ) : (
            <div className="rounded-[24px] border border-dashed border-slate-300 bg-slate-50 px-5 py-8 text-sm text-slate-500">
              Комментариев пока нет. Начните обсуждение первым.
            </div>
          )}
        </Card>

        <Card className="space-y-5">
          <div>
            <h2 className="text-xl font-semibold text-ink">Новый комментарий</h2>
            <p className="mt-1 text-sm text-slate-500">Короткие, понятные комментарии читаются лучше и не разваливают обсуждение.</p>
          </div>
          <form className="space-y-4" noValidate onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
            <Field label="Комментарий" error={form.formState.errors.content?.message}>
              <Textarea {...form.register("content")} className="min-h-[200px]" placeholder="Поделитесь мнением, уточните детали или предложите решение." />
            </Field>
            {errorMessage ? <div className="text-sm text-red-600">{errorMessage}</div> : null}
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Отправляем..." : "Добавить комментарий"}
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
}
