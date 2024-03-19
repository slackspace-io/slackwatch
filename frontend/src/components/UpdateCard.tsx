import React from 'react';
import { Card } from '@/components/ui/card';
import { unstable_noStore as noStore } from 'next/cache';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import RefreshButton from "@/components/RefreshButton";


interface UpdateCardProps {
  update: {
    name: string,
    exclude_pattern: string,
    git_ops_repo: string,
    include_pattern: string,
    update_available: string,
    image: string,
    last_scanned: string,
    namespace: string,
    current_version: string,
    latest_version: string,
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
    const repo = data.get('gitopsRepo');
    //create url
    //fetch url
    const params = new URLSearchParams({
        name: name as string,
        image: image as string,
        newTag: newTag as string,
        currentTag: currentTag as string,
        namespace: namespace as string,
        repo: repo as string
    });
    const response = await fetch(`${baseUrl}/api/container/update?${params}`);
    if (!response.ok) {
        throw new Error('Failed to fetch data');
        return false;
    }
    return true;
    //notify success

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
            {update.name} has an update available. Click the button to update.
        </AlertDescription>
    </Alert>
);

const UpdateCard: React.FC<UpdateCardProps> = ({ update }) => (
  <Card className={`mb-4 p-4 shadow-lg rounded-lg ${update.update_available ? 'border-l-4 border-green-500' : 'border'}`}>
    <div className="flex justify-between items-center">
      <p className="text-lg font-bold">{update.name}</p>
      {update.update_available && update.latest_version && (
          <form action={handleUpdate}>
        <input name="containerName" type="hidden" value={update.name} />
        <input name="newTag" type="hidden" value={update.latest_version} />
        <input name="image" type="hidden" value={update.image} />
        <input name="currentTag" type="hidden" value={update.current_version} />
        <input name="namespace" type="hidden" value={update.namespace} />
        <input name="gitopsRepo" type="hidden" value={update.git_ops_repo} />
        <button className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition duration-150 ease-in-out">
          Update
        </button>
          </form>
      )}
    </div>
      <div className="flex justify-center items-end">
          <RefreshButton />
      </div>
    <div className="mt-4">
      {update.current_version && <p>Current Tag: <span className="font-semibold">{update.current_version}</span></p>}
      {update.latest_version && <p className="text-green-500">New Tag: <span className="font-semibold">{update.latest_version}</span></p>}
      <p>Image: <span className="font-semibold">{update.image}</span></p>
      <p>Pod Name: <span className="font-semibold">{update.name}</span></p>
      <p>Time Scanned: <span className="font-semibold">{update.last_scanned}</span></p>
        <p>Namespace: <span className="font-semibold">{update.namespace}</span></p>
      {update.last_scanned && <p>Notification Sent At: <span className="font-semibold">{update.last_scanned}</span></p>}
    </div>
  </Card>

);

export default UpdateCard;
