import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useId, useRef, useState } from "react";

import { useSession } from "@/app/providers";
import type { ChatMessage } from "@/entities/chat/model";
import { api, toAssetUrl } from "@/shared/api/client";
import { classNames, formatDate } from "@/shared/lib/format";
import { Button, Card, Input, PageHeader, Textarea } from "@/shared/ui/primitives";

export function ChatPage() {
  const { user, selectedHouseId } = useSession();
  const queryClient = useQueryClient();
  const inputId = useId();
  const [content, setContent] = useState("");
  const [image, setImage] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState("");
  const listRef = useRef<HTMLDivElement | null>(null);

  const chatQuery = useQuery({
    queryKey: ["chat", selectedHouseId],
    queryFn: async () => {
      const response = await api.get<{ items: ChatMessage[] }>(`/houses/${selectedHouseId}/chat/messages`);
      return response.items;
    },
    enabled: Boolean(selectedHouseId),
    refetchInterval: 3000,
  });

  useEffect(() => {
    const element = listRef.current;
    if (!element) {
      return;
    }
    element.scrollTop = element.scrollHeight;
  }, [chatQuery.data]);

  useEffect(() => {
    if (!image) {
      setPreviewUrl("");
      return;
    }
    const objectUrl = URL.createObjectURL(image);
    setPreviewUrl(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [image]);

  const mutation = useMutation({
    mutationFn: async () => {
      const payload = new FormData();
      payload.append("content", content);
      if (image) {
        payload.append("image", image);
      }
      return api.post(`/houses/${selectedHouseId}/chat/messages`, payload);
    },
    onSuccess() {
      setContent("");
      setImage(null);
      queryClient.invalidateQueries({ queryKey: ["chat", selectedHouseId] });
    },
  });

  const canSend = content.trim().length > 0 || Boolean(image);
  const errorMessage = mutation.error instanceof Error ? mutation.error.message : "";

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Чат дома"
        title="Общий чат соседей"
        description="Внутри этого дома все участники видят общую переписку, изображения и имена авторов. Лента сообщений обновляется автоматически."
      />

      <Card className="grid gap-6 p-0 lg:grid-cols-[1.25fr_0.75fr]">
        <div className="border-b border-slate-100 p-5 lg:border-b-0 lg:border-r">
          <div className="mb-4 flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-ink">Сообщения</h2>
              <p className="text-sm text-slate-500">Живой разговор соседей в формате мессенджера.</p>
            </div>
            <span className="rounded-full bg-sand px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-ink">онлайн-лента</span>
          </div>

          <div ref={listRef} className="flex max-h-[68vh] flex-col gap-4 overflow-y-auto pr-2">
            {chatQuery.data?.length ? (
              chatQuery.data.map((message) => {
                const ownMessage = message.authorId === user?.id;

                return (
                  <div key={message.id} className={classNames("flex", ownMessage ? "justify-end" : "justify-start")}>
                    <article
                      className={classNames(
                        "max-w-[88%] rounded-[26px] px-4 py-3 shadow-soft md:max-w-[72%]",
                        ownMessage ? "bg-ink text-white" : "bg-slate-100 text-ink",
                      )}
                    >
                      <div className={classNames("mb-2 flex items-center gap-2 text-xs", ownMessage ? "text-sand/80" : "text-slate-500")}>
                        <span className="font-semibold">{message.authorName}</span>
                        <span>{formatDate(message.createdAt)}</span>
                      </div>
                      {message.content ? <p className="whitespace-pre-wrap text-sm leading-6">{message.content}</p> : null}
                      {message.imageUrl ? (
                        <div className={classNames("mt-3 overflow-hidden rounded-[20px]", ownMessage ? "bg-white/10" : "bg-white")}>
                          <img src={toAssetUrl(message.imageUrl)} alt="Изображение из чата" className="max-h-80 w-full object-cover" />
                        </div>
                      ) : null}
                    </article>
                  </div>
                );
              })
            ) : (
              <div className="rounded-[24px] border border-dashed border-slate-300 bg-slate-50 px-5 py-8 text-center text-sm text-slate-500">
                Чат пока пуст. Отправьте первое сообщение и задайте тон обсуждению.
              </div>
            )}
          </div>
        </div>

        <div className="space-y-5 p-5">
          <div className="space-y-2">
            <h2 className="text-xl font-semibold text-ink">Написать в чат</h2>
            <p className="text-sm text-slate-500">Можно отправить текст, картинку или сразу оба варианта в одном сообщении.</p>
          </div>

          <form
            className="space-y-4"
            onSubmit={(event) => {
              event.preventDefault();
              if (!canSend || mutation.isPending) {
                return;
              }
              mutation.mutate();
            }}
          >
            <Textarea
              value={content}
              onChange={(event) => setContent(event.target.value)}
              className="min-h-[180px]"
              placeholder="Например: соседи, завтра в 19:00 обсуждаем освещение у подъезда."
            />

            <div className="flex items-center gap-3">
              <input
                id={inputId}
                type="file"
                accept="image/png,image/jpeg,image/webp,image/gif"
                className="hidden"
                onChange={(event) => setImage(event.target.files?.[0] ?? null)}
              />
              <label
                htmlFor={inputId}
                className="inline-flex h-12 w-12 cursor-pointer items-center justify-center rounded-2xl border border-slate-200 bg-white text-ink transition hover:border-moss hover:text-moss"
                title="Прикрепить изображение"
              >
                <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.9" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21.44 11.05 12 20.5a6 6 0 0 1-8.49-8.49l10.6-10.61a4 4 0 0 1 5.66 5.66L9.17 17.66a2 2 0 0 1-2.83-2.83l9.2-9.19" />
                </svg>
              </label>
              <div className="text-sm text-slate-500">{image ? image.name : "Прикрепить изображение"}</div>
            </div>

            {previewUrl ? (
              <div className="overflow-hidden rounded-[24px] border border-slate-200 bg-slate-50">
                <img src={previewUrl} alt="Предпросмотр сообщения" className="h-56 w-full object-cover" />
              </div>
            ) : null}
            {errorMessage ? <div className="text-sm text-red-600">{errorMessage}</div> : null}
            <Button type="submit" disabled={!canSend || mutation.isPending} className="w-full">
              {mutation.isPending ? "Отправляем..." : "Отправить сообщение"}
            </Button>
          </form>
        </div>
      </Card>
    </div>
  );
}
