<script>
  import NdaDocument from '$lib/NdaDocument.svelte';
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

  let downloading = $state(false);

  async function handleDownload() {
    downloading = true;
    try {
      await generatePdf('nda-document');
    } finally {
      downloading = false;
    }
  }
</script>

<div class="flex flex-col h-screen bg-gray-50">
  <header class="flex items-center justify-between px-6 py-4 bg-white border-b shadow-sm shrink-0">
    <div>
      <h1 class="text-xl font-semibold text-gray-900">Mutual NDA Creator</h1>
      <p class="text-sm text-gray-500">Fill in the details to generate your Mutual Non-Disclosure Agreement</p>
    </div>
    <button
      onclick={handleDownload}
      disabled={downloading}
      class="px-5 py-2 bg-gray-900 text-white text-sm font-medium rounded-md hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
    >
      {downloading ? 'Generating...' : 'Download PDF'}
    </button>
  </header>

  <div class="flex flex-1 min-h-0">
    <!-- Form Panel -->
    <div class="w-2/5 overflow-y-auto border-r bg-white p-6 space-y-6">

      <!-- Party 1 -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Party 1</h2>
        <div class="space-y-3">
          <div>
            <label for="p1-name" class="block text-xs text-gray-500 mb-1">Name</label>
            <input id="p1-name" type="text" bind:value={form.party1.name} placeholder="Full name" class="input" />
          </div>
          <div>
            <label for="p1-title" class="block text-xs text-gray-500 mb-1">Title</label>
            <input id="p1-title" type="text" bind:value={form.party1.title} placeholder="e.g. CEO" class="input" />
          </div>
          <div>
            <label for="p1-company" class="block text-xs text-gray-500 mb-1">Company</label>
            <input id="p1-company" type="text" bind:value={form.party1.company} placeholder="Company name" class="input" />
          </div>
          <div>
            <label for="p1-address" class="block text-xs text-gray-500 mb-1">Notice Address</label>
            <textarea id="p1-address" bind:value={form.party1.address} placeholder="Mailing address" rows="2" class="input resize-none"></textarea>
          </div>
        </div>
      </section>

      <hr class="border-gray-100" />

      <!-- Party 2 -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Party 2</h2>
        <div class="space-y-3">
          <div>
            <label for="p2-name" class="block text-xs text-gray-500 mb-1">Name</label>
            <input id="p2-name" type="text" bind:value={form.party2.name} placeholder="Full name" class="input" />
          </div>
          <div>
            <label for="p2-title" class="block text-xs text-gray-500 mb-1">Title</label>
            <input id="p2-title" type="text" bind:value={form.party2.title} placeholder="e.g. CEO" class="input" />
          </div>
          <div>
            <label for="p2-company" class="block text-xs text-gray-500 mb-1">Company</label>
            <input id="p2-company" type="text" bind:value={form.party2.company} placeholder="Company name" class="input" />
          </div>
          <div>
            <label for="p2-address" class="block text-xs text-gray-500 mb-1">Notice Address</label>
            <textarea id="p2-address" bind:value={form.party2.address} placeholder="Mailing address" rows="2" class="input resize-none"></textarea>
          </div>
        </div>
      </section>

      <hr class="border-gray-100" />

      <!-- Agreement Details -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Agreement Details</h2>
        <div class="space-y-3">
          <div>
            <label for="purpose" class="block text-xs text-gray-500 mb-1">Purpose</label>
            <textarea id="purpose" bind:value={form.purpose} placeholder="How Confidential Information may be used" rows="2" class="input resize-none"></textarea>
          </div>
          <div>
            <label for="effective-date" class="block text-xs text-gray-500 mb-1">Effective Date</label>
            <input id="effective-date" type="date" bind:value={form.effectiveDate} class="input" />
          </div>
        </div>
      </section>

      <hr class="border-gray-100" />

      <!-- MNDA Term -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">MNDA Term</h2>
        <div class="space-y-2">
          <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
            <input type="radio" bind:group={form.mndaTermType} value="expires" class="accent-gray-900" />
            Expires after
            <input
              id="mnda-term-years"
              type="number"
              bind:value={form.mndaTermYears}
              min="1"
              aria-label="MNDA term years"
              class="input w-16 text-center py-0.5"
              disabled={form.mndaTermType !== 'expires'}
            />
            year(s)
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
            <input type="radio" bind:group={form.mndaTermType} value="continues" class="accent-gray-900" />
            Continues until terminated
          </label>
        </div>
      </section>

      <hr class="border-gray-100" />

      <!-- Term of Confidentiality -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Term of Confidentiality</h2>
        <div class="space-y-2">
          <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
            <input type="radio" bind:group={form.confidentialityTermType} value="years" class="accent-gray-900" />
            <input
              type="number"
              bind:value={form.confidentialityTermYears}
              min="1"
              aria-label="Confidentiality term years"
              class="input w-16 text-center py-0.5"
              disabled={form.confidentialityTermType !== 'years'}
            />
            year(s) from Effective Date
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
            <input type="radio" bind:group={form.confidentialityTermType} value="perpetuity" class="accent-gray-900" />
            In perpetuity
          </label>
        </div>
      </section>

      <hr class="border-gray-100" />

      <!-- Legal -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wide mb-3">Legal</h2>
        <div class="space-y-3">
          <div>
            <label for="governing-law" class="block text-xs text-gray-500 mb-1">Governing Law (State)</label>
            <input id="governing-law" type="text" bind:value={form.governingLaw} placeholder="e.g. California" class="input" />
          </div>
          <div>
            <label for="jurisdiction" class="block text-xs text-gray-500 mb-1">Jurisdiction</label>
            <input id="jurisdiction" type="text" bind:value={form.jurisdiction} placeholder="e.g. San Francisco, California" class="input" />
          </div>
          <div>
            <label for="modifications" class="block text-xs text-gray-500 mb-1">
              Modifications <span class="text-gray-400">(optional)</span>
            </label>
            <textarea id="modifications" bind:value={form.modifications} placeholder="Any modifications to the Standard Terms" rows="2" class="input resize-none"></textarea>
          </div>
        </div>
      </section>

    </div>

    <!-- Preview Panel -->
    <div class="flex-1 overflow-y-auto bg-gray-100 p-6">
      <div id="nda-document" class="shadow-md">
        <NdaDocument {form} />
      </div>
    </div>
  </div>
</div>

<style>
  :global(.input) {
    width: 100%;
    padding: 0.375rem 0.625rem;
    border: 1px solid #d1d5db;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    color: #111827;
    background: #fff;
    outline: none;
    transition: border-color 0.15s;
  }

  :global(.input:focus) {
    border-color: #6b7280;
    box-shadow: 0 0 0 2px rgba(107, 114, 128, 0.2);
  }

  :global(.input:disabled) {
    background: #f9fafb;
    color: #9ca3af;
  }
</style>
