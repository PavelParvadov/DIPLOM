import { Link } from "react-router-dom";

import type { Post } from "@/entities/post/model";
import { toAssetUrl } from "@/shared/api/client";
import { classNames, formatDate } from "@/shared/lib/format";
import { SecondaryButton, Card } from "@/shared/ui/primitives";

interface PostCardProps {
  post: Post;
  immersive?: boolean;
}

export function PostCard({ post, immersive = false }: PostCardProps) {
  const hasImage = Boolean(post.imageUrl);
  const isImmersive = immersive && hasImage;

  return (
    <Card
      className={classNames(
        "space-y-4 overflow-hidden",
        isImmersive ? "snap-start rounded-[32px] p-0" : "rounded-[28px]",
      )}
    >
      {hasImage ? (
        <img
          src={toAssetUrl(post.imageUrl)}
          alt={post.title}
          className={classNames("w-full object-cover", isImmersive ? "h-[42vh] md:h-[54vh]" : "h-60")}
        />
      ) : null}

      <div className={classNames("space-y-4", isImmersive ? "p-6 md:p-8" : "")}>
        <div className="flex flex-wrap items-center gap-3 text-xs uppercase tracking-[0.2em] text-slate-500">
          <span>{post.categoryName}</span>
          <span>{formatDate(post.createdAt)}</span>
        </div>
        <div className="space-y-3">
          <Link to={`/posts/${post.id}`} className="block">
            <h3 className={classNames("font-semibold text-ink transition hover:text-moss", isImmersive ? "text-2xl md:text-3xl" : "text-xl")}>
              {post.title}
            </h3>
          </Link>
          <p className={classNames("leading-7 text-slate-600", isImmersive ? "text-base" : "line-clamp-4 text-sm leading-6")}>{post.content}</p>
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="text-sm font-medium text-slate-500">Автор: {post.authorName}</div>
          <Link to={`/posts/${post.id}`}>
            <SecondaryButton type="button">
              Комментарии{post.commentsCount ? ` (${post.commentsCount})` : ""}
            </SecondaryButton>
          </Link>
        </div>
      </div>
    </Card>
  );
}
