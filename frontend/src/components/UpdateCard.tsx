import React from 'react';
import { Card } from '@/components/ui/card';
import { unstable_noStore as noStore } from 'next/cache';
import {Alert, AlertDescription, AlertTitle} from '@/components/ui/alert';

interface UpdateCardProps {
  update: {
    containerName: string;
    currentTag?: string;
    newTag?: string;
    image: string;
    podName: string;
    timeScanned: string;
    namespace: string;
    updateAvailable: boolean;
    sentTime?: string;
    repo: string;
  };
}

async function handleUpdate(data: FormData) {
    "use server";
    noStore();
    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
    if (!baseUrl) {
        console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
        return [];
    }
    const name = data.get('containerName');
    const image = data.get('image');
    const newTag = data.get('newTag');
    const currentTag = data.get('currentTag');
    const namespace = data.get('namespace');
    const repo = data.get('repo');
    //create url
    //fetch url
    const params = new URLSearchParams({
        containerName: name as string,
        image: image as string,
        newTag: newTag as string,
        currentTag: currentTag as string,
        namespace: namespace as string,
        repo: repo as string
    });
    const response = await fetch(`${baseUrl}/api/container/update?${params}`);
    if (!response.ok) {
        throw new Error('Failed to fetch data');
    }
    console.log("hi")

}

//const handleUpdate = async (data: FormData) => {
//    "use server";
//    console.log('Update clicked');
//    console.log(data)
//    //call PostUpdateRequest
//    const response = await PostUpdateRequest(data);
//    //log data
//
//}
const AlertMessage: React.FC<UpdateCardProps> = ({ update }) => (
    <Alert>
        <AlertTitle>Update Available</AlertTitle>
        <AlertDescription>
            {update.containerName} has an update available. Click the button to update.
        </AlertDescription>
    </Alert>
);

const UpdateCard: React.FC<UpdateCardProps> = ({ update }) => (
  <Card className={`mb-4 p-4 shadow-lg rounded-lg ${update.updateAvailable ? 'border-l-4 border-green-500' : 'border'}`}>
    <div className="flex justify-between items-center">
      <p className="text-lg font-bold">{update.containerName}</p>
      {update.updateAvailable && update.newTag && (
          <form action={handleUpdate}>
        <input name="containerName" type="hidden" value={update.containerName} />
        <input name="newTag" type="hidden" value={update.newTag} />
        <input name="image" type="hidden" value={update.image} />
        <input name="currentTag" type="hidden" value={update.currentTag} />
        <input name="namespace" type="hidden" value={update.namespace} />
        <button className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition duration-150 ease-in-out">
          Update
        </button>
          </form>
      )}
    </div>
    <div className="mt-4">
      {update.currentTag && <p>Current Tag: <span className="font-semibold">{update.currentTag}</span></p>}
      {update.newTag && <p className="text-green-500">New Tag: <span className="font-semibold">{update.newTag}</span></p>}
      <p>Image: <span className="font-semibold">{update.image}</span></p>
      <p>Pod Name: <span className="font-semibold">{update.podName}</span></p>
      <p>Time Scanned: <span className="font-semibold">{update.timeScanned}</span></p>
      {update.sentTime && <p>Notification Sent At: <span className="font-semibold">{update.sentTime}</span></p>}
    </div>
  </Card>
);

export default UpdateCard;
