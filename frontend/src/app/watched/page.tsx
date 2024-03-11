import { Card } from '@/components/ui/card';
import Link from 'next/link'; // Import Link from next/link
import { unstable_noStore as noStore } from 'next/cache'



interface Container {
    containerName: string;
    excludePattern: string;
    image: string;
    includePattern: string;
    podName: string;
    timeScanned: string;
}

async function getContainers(): Promise<Container[]> {
    noStore();
    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
    console.log("found base url get containers", baseUrl);
    // Check if baseUrl is defined
    if (!baseUrl) {
        console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
        return []; // Return empty array or mock data
    }
    const res = await fetch(`${baseUrl}/api/containers`);
    if (!res.ok) {
        throw new Error('Failed to fetch data');
    }
    return res.json();
}

export default async function Watched() {
    const containers = await getContainers();
    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
    console.log("found base url in watched", baseUrl);
    return (
        <main className="p-4">
            <div className="mb-4">
                <Link href="/" className="text-blue-500 hover:text-blue-700">
                    Available Updates
                </Link>
            </div>
            {containers.map((container: Container, index: number) => (
                <Card key={index} className="mb-4 p-4 shadow-lg">
                    <p className="text-lg font-bold">
                        Container Name: {container.containerName}
                    </p>
                    <p>Image: {container.image}</p>
                    <p>Pod Name: {container.podName}</p>
                    <p>Time Scanned: {container.timeScanned}</p>
                    <p>Include Pattern: {container.includePattern}</p>
                    <p>Exclude Pattern: {container.excludePattern}</p>
                </Card>
            ))}
        </main>
    );
}