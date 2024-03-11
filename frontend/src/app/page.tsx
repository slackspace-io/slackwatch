import { Card } from '@/components/ui/card';
import Link from 'next/link'; // Import Link from next/link

interface Update {
  containerName: string;
  currentTag: string;
  newTag: string;
  foundAt: string;
}

async function getData(): Promise<Update[]> {
  // Use the environment variable
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  // Check if baseUrl is defined
  if (!baseUrl) {
    console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
    return []; // Return empty array or mock data
  }
  const res = await fetch(`${baseUrl}/api/imageUpdates`);
  if (!res.ok) {
    throw new Error('Failed to fetch data');
  }
  return res.json();
}

export default async function Page() {
  const data = await getData();

  return (
    <main className="p-4">
      <div className="mb-4">
        {/* Link to the Dashboard page */}
        <Link href="/watched" className="text-blue-500 hover:text-blue-700">
          All Watched
        </Link>
      </div>
      {data.map((update: Update, index: number) => (
        <Card key={index} className="mb-4 p-4 shadow-lg">
          <p className="text-lg font-bold">
            Container: {update.containerName}
          </p>
          <p>Current Tag: {update.currentTag}</p>
          <p>New Tag: {update.newTag}</p>
          <p>Found At: {update.foundAt}</p>
        </Card>
      ))}
    </main>
  );
}