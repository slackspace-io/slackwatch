import React from 'react';
import { Card } from '@/components/ui/card';

interface UpdateCardProps {
  update: {
    containerName: string;
    currentTag?: string;
    newTag?: string;
    image: string;
    podName: string;
    timeScanned: string;
    updateAvailable: boolean;
    sentTime?: string;
  };
}

const UpdateCard: React.FC<UpdateCardProps> = ({ update }) => (
  <Card className={`mb-4 p-4 shadow-lg rounded-lg ${update.updateAvailable ? 'border-l-4 border-green-500' : 'border'}`}>
    <div className="flex justify-between items-center">
      <p className="text-lg font-bold">{update.containerName}</p>
      {update.updateAvailable && update.newTag && (
        <button className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition duration-150 ease-in-out">
          Update
        </button>
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