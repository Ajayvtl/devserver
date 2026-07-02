# DevServer Web

This is the Next.js UI for DevServer.

## Structure

- `app/`: routed pages and shared layout
- `components/`: reusable UI, shell components, and flow screens
- `lib/`: types, mock data, and service facades
- `hooks/`: client-side hooks for future API replacement

## Flow

- `/`: bootstrap router with loading animation and redirect decision
- `/setup`: first-time setup wizard
- `/login`: sign-in flow
- `/dashboard`: main product shell
- route-level `loading.tsx` files provide polished transitional states

## Scripts

- `npm run dev`
- `npm run build`
- `npm run start`
