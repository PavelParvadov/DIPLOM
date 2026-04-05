import { useQuery } from "@tanstack/react-query";

import type { House } from "@/entities/house/model";
import { api } from "@/shared/api/client";

export async function fetchHouses() {
  const response = await api.get<{ items: House[] }>("/houses");
  return response.items;
}

export function useHouses(enabled = true) {
  return useQuery({
    queryKey: ["houses"],
    queryFn: fetchHouses,
    enabled,
  });
}
