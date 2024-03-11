import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

export async function middleware(req: NextRequest) {
  const apiUrl = process.env.NEXT_PUBLIC_API_BASE_URL + 'api/imageUpdates';
  const response = await fetch(apiUrl);
  const data = await response.json();

  return NextResponse.json(data);
}