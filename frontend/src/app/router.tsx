import { createBrowserRouter, Navigate, Outlet, RouterProvider } from "react-router-dom";

import { AppLayout } from "@/app/layout";
import { useSession } from "@/app/providers";
import { useHouses } from "@/features/houses/queries";
import { AdminCategoriesPage } from "@/pages/admin-categories-page";
import { AdminInvitesPage } from "@/pages/admin-invites-page";
import { AdminPostsPage } from "@/pages/admin-posts-page";
import { ChatPage } from "@/pages/chat-page";
import { FeedPage } from "@/pages/feed-page";
import { JoinHousePage } from "@/pages/join-house-page";
import { LoginPage } from "@/pages/login-page";
import { NewPostPage } from "@/pages/new-post-page";
import { PostDetailPage } from "@/pages/post-detail-page";
import { RegisterPage } from "@/pages/register-page";

function PublicOnly() {
  const { user } = useSession();
  return user ? <Navigate to="/posts" replace /> : <Outlet />;
}

function Protected() {
  const { user } = useSession();
  return user ? <AppLayout /> : <Navigate to="/login" replace />;
}

function AdminOnly() {
  const { user, selectedHouseId } = useSession();
  const housesQuery = useHouses(Boolean(user));

  if (!user) {
    return <Navigate to="/login" replace />;
  }
  if (housesQuery.isLoading) {
    return null;
  }

  const selectedHouse = housesQuery.data?.find((house) => house.id === selectedHouseId);
  if (!selectedHouse || selectedHouse.role !== "admin") {
    return <Navigate to="/posts" replace />;
  }

  return <Outlet />;
}

function RootRedirect() {
  const { user } = useSession();
  return <Navigate to={user ? "/posts" : "/login"} replace />;
}

export function AppRouter() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <RootRedirect />,
    },
    {
      element: <PublicOnly />,
      children: [
        { path: "/login", element: <LoginPage /> },
        { path: "/register", element: <RegisterPage /> },
      ],
    },
    {
      element: <Protected />,
      children: [
        { path: "/join", element: <JoinHousePage /> },
        { path: "/posts", element: <FeedPage /> },
        { path: "/chat", element: <ChatPage /> },
        { path: "/posts/new", element: <NewPostPage /> },
        { path: "/posts/:postId", element: <PostDetailPage /> },
        {
          element: <AdminOnly />,
          children: [
            { path: "/admin/categories", element: <AdminCategoriesPage /> },
            { path: "/admin/posts", element: <AdminPostsPage /> },
            { path: "/admin/invites", element: <AdminInvitesPage /> },
          ],
        },
      ],
    },
  ]);

  return <RouterProvider router={router} />;
}
