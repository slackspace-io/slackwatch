import { unstable_noStore as noStore } from 'next/cache';
import UpdateCard from '@/components/UpdateCard';
import RefreshButton from "@/components/RefreshButton";
import React from "react";


interface CombinedData {
    name: string,
    exclude_pattern: string,
    git_ops_repo: string,
    include_pattern: string,
    update_available: string,
    image: string,
    last_scanned: string,
    git_directory: string,
    namespace: string,
    current_version: string,
    latest_version: string,
}
async function getData(): Promise<CombinedData[]> {
  noStore();
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!baseUrl) {
    console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
    return [];
  }
  const res = await fetch(`${baseUrl}/api/workloads/all`);
  if (!res.ok) {
    throw new Error('Failed to fetch data');
  }
  return res.json();
}

async function refreshData() {
    'use server';
  noStore();
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!baseUrl) {
    console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
    return [];
  }
  const res = await fetch(`${baseUrl}/api/workloads/refresh`);
  if (!res.ok) {
    throw new Error('Failed to fetch data');
  }
  return res.json();
}

export default async function Page() {
  let data = await getData();

//  data = data.sort((a, b) => Number(b.update_available) - Number(a.update_available));
    //order by two strings which goes first
    data = data.sort((a, b) => { return a.update_available.localeCompare(b.update_available) });
  return (
      <main className="p-4">
          <div className="flex justify-center items-end">
              <form action={refreshData} >
                  <RefreshButton/>

              </form>
          </div>
          {data.map((update, index) => (
              <UpdateCard key={index} update={update}/>
          ))}
      </main>
  );
}
