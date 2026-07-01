<script>
  import NdaDocument from '$lib/NdaDocument.svelte';
  import Chat from '$lib/Chat.svelte';
  import { generatePdf } from '$lib/generatePdf.js';

  let form = $state({
    party1: { name: '', title: '', company: '', address: '' },
    party2: { name: '', title: '', company: '', address: '' },
    effectiveDate: '',
    mndaTermType: 'expires',
    mndaTermYears: 1,
    confidentialityTermType: 'years',
    confidentialityTermYears: 3,
    purpose: '',
    governingLaw: '',
    jurisdiction: '',
    modifications: '',
  });

  function handleDownload() {
    generatePdf();
  }
</script>

<div class="app-root flex flex-col h-screen bg-gray-50">
  <header class="app-header flex items-center justify-between px-6 py-4 bg-white border-b shadow-sm shrink-0">
    <div>
      <h1 class="text-xl font-semibold text-gray-900">Mutual NDA Creator</h1>
      <p class="text-sm text-gray-500">Fill in the details to generate your Mutual Non-Disclosure Agreement</p>
    </div>
    <button
      onclick={handleDownload}
      class="px-5 py-2 bg-gray-900 text-white text-sm font-medium rounded-md hover:bg-gray-700 transition-colors"
    >
      Save as PDF
    </button>
  </header>

  <div class="flex flex-1 min-h-0">
    <!-- Chat Panel -->
    <div class="chat-panel w-2/5 min-h-0 overflow-hidden border-r bg-white">
      <Chat onFormData={(data) => (form = data)} />
    </div>

    <!-- Preview Panel -->
    <div class="preview-panel flex-1 overflow-y-auto bg-gray-100 p-6">
      <div id="nda-document" class="shadow-md">
        <NdaDocument {form} />
      </div>
    </div>
  </div>
</div>

<style>
  @media print {
    :global(.app-header),
    :global(.chat-panel) {
      display: none !important;
    }

    :global(.app-root),
    :global(.app-root > div) {
      display: block !important;
      height: auto !important;
      overflow: visible !important;
      background: white !important;
    }

    :global(#nda-document) {
      box-shadow: none !important;
      max-width: 100% !important;
    }

    :global(.preview-panel) {
      overflow: visible !important;
      background: white !important;
      padding: 0 !important;
    }
  }
</style>
