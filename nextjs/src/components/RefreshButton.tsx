import { CombinedData } from '@/types/workloads';
import {unstable_noStore as noStore} from "next/dist/server/web/spec-extension/unstable-no-store";

async function refreshSingleData(update: CombinedData) {
    'use server';
    noStore();
    console.log('refreshSingleData');
    const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
    if (!baseUrl) {
        console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
        return [];
    }
    const params = new URLSearchParams({
        name: update.name,
        exclude_pattern: update.exclude_pattern,
        git_ops_repo: update.git_ops_repo,
        include_pattern: update.include_pattern,
        update_available: update.update_available,
        image: update.image,
        last_scanned: update.last_scanned,
        git_directory: update.git_directory,
        namespace: update.namespace,
        current_version: update.current_version,
        latest_version: update.latest_version,
    });
    console.log(params.toString());
    const res = await fetch(`${baseUrl}/api/update/workload?${params}`);
    if (!res.ok) {
        throw new Error('Failed to fetch data');
    }
    return res.json();

}

const RefreshButton = ({update}: {update: CombinedData}) => {
  const handleRefresh = async () => {
      'use server';
    console.log('Refresh clicked');
    console.log(update);
    await refreshSingleData(update);
  }
    return (
        <div>
        <form action={handleRefresh}>
            <button>Refresh</button>
        </form>
        </div>
    );
}
export default RefreshButton;
