import {clsx, type ClassValue} from "clsx";
import {twMerge} from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function extractInitials(firstName: string, lastName: string): string {
  return `${firstName.charAt(0)} ${lastName.charAt(0)}`.toUpperCase();
}

export function formatDate(input: string): string {
  const date = new Date(input);

  return date.toLocaleDateString("en-US", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}
