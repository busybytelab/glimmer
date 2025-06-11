<script lang="ts">
    import { sessionImportExportService } from '$lib/services/sessionImportExport';
    import FormButton from '$components/common/FormButton.svelte';
    import { toast } from '$lib/stores/toast';

    export let sessionId: string;
    export let disabled = false;
    let isLoading = false;

    async function handleExport() {
        isLoading = true;
        try {
            const exportData = await sessionImportExportService.exportPracticeSession(sessionId);
            sessionImportExportService.downloadExportedSession(exportData);
            toast.success('Session exported successfully');
        } catch (error) {
            console.error('Failed to export session:', error);
            toast.error('Failed to export session');
        } finally {
            isLoading = false;
        }
    }
</script>

<FormButton 
    variant="secondary" 
    type="button"
    on:click={handleExport} 
    {disabled}
    {isLoading}
    loadingText="Exporting..."
>
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
    </svg>
    Export
</FormButton> 