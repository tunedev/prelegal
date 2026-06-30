<script>
  let { form } = $props();

  function formatDate(dateStr) {
    if (!dateStr) return '[Effective Date not set]';
    const [y, m, d] = dateStr.split('-').map(Number);
    return new Date(y, m - 1, d).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }

  let effectiveDate = $derived(formatDate(form.effectiveDate));

  let mndaTerm = $derived(
    form.mndaTermType === 'expires'
      ? `${form.mndaTermYears} year(s) from Effective Date`
      : 'until terminated in accordance with the terms of the MNDA',
  );

  let confidentialityTerm = $derived(
    form.confidentialityTermType === 'years'
      ? `${form.confidentialityTermYears} year(s) from Effective Date, but in the case of trade secrets until Confidential Information is no longer considered a trade secret under applicable laws`
      : 'in perpetuity',
  );

  let purpose = $derived(form.purpose || '[Purpose not specified]');
  let governingLaw = $derived(form.governingLaw || '[Governing Law not specified]');
  let jurisdiction = $derived(form.jurisdiction || '[Jurisdiction not specified]');
</script>

<div class="document">
  <!-- ═══════════════════════════════ COVER PAGE ═══════════════════════════════ -->
  <h1>Mutual Non-Disclosure Agreement</h1>

  <section class="intro">
    <h2>USING THIS MUTUAL NON-DISCLOSURE AGREEMENT</h2>
    <p>
      This Mutual Non-Disclosure Agreement (the "MNDA") consists of: (1) this Cover Page
      ("Cover Page") and (2) the Common Paper Mutual NDA Standard Terms Version 1.0
      ("Standard Terms") identical to those posted at commonpaper.com/standards/mutual-nda/1.0.
      Any modifications of the Standard Terms should be made on the Cover Page, which will
      control over conflicts with the Standard Terms.
    </p>
  </section>

  <section class="cover-field">
    <h3>Purpose</h3>
    <p class="field-label">How Confidential Information may be used</p>
    <p>{purpose}</p>
  </section>

  <section class="cover-field">
    <h3>Effective Date</h3>
    <p>{effectiveDate}</p>
  </section>

  <section class="cover-field">
    <h3>MNDA Term</h3>
    <p class="field-label">The length of this MNDA</p>
    <p>
      {#if form.mndaTermType === 'expires'}
        <span class="checkbox checked">&#9746;</span> Expires {form.mndaTermYears} year(s) from Effective Date.<br />
        <span class="checkbox">&#9744;</span> Continues until terminated in accordance with the terms of the MNDA.
      {:else}
        <span class="checkbox">&#9744;</span> Expires {form.mndaTermYears} year(s) from Effective Date.<br />
        <span class="checkbox checked">&#9746;</span> Continues until terminated in accordance with the terms of the MNDA.
      {/if}
    </p>
  </section>

  <section class="cover-field">
    <h3>Term of Confidentiality</h3>
    <p class="field-label">How long Confidential Information is protected</p>
    <p>
      {#if form.confidentialityTermType === 'years'}
        <span class="checkbox checked">&#9746;</span> {form.confidentialityTermYears} year(s) from Effective Date, but in the case of trade secrets until Confidential Information is no longer considered a trade secret under applicable laws.<br />
        <span class="checkbox">&#9744;</span> In perpetuity.
      {:else}
        <span class="checkbox">&#9744;</span> {form.confidentialityTermYears} year(s) from Effective Date, but in the case of trade secrets until Confidential Information is no longer considered a trade secret under applicable laws.<br />
        <span class="checkbox checked">&#9746;</span> In perpetuity.
      {/if}
    </p>
  </section>

  <section class="cover-field">
    <h3>Governing Law &amp; Jurisdiction</h3>
    <p>Governing Law: {governingLaw}</p>
    <p>Jurisdiction: {jurisdiction}</p>
  </section>

  <section class="cover-field">
    <h3>MNDA Modifications</h3>
    <p>{form.modifications || 'None'}</p>
  </section>

  <p class="signing-intro">
    By signing this Cover Page, each party agrees to enter into this MNDA as of the Effective Date.
  </p>

  <table class="signature-table">
    <thead>
      <tr>
        <th></th>
        <th>PARTY 1</th>
        <th>PARTY 2</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="row-label">Signature</td>
        <td></td>
        <td></td>
      </tr>
      <tr>
        <td class="row-label">Print Name</td>
        <td>{form.party1.name || ''}</td>
        <td>{form.party2.name || ''}</td>
      </tr>
      <tr>
        <td class="row-label">Title</td>
        <td>{form.party1.title || ''}</td>
        <td>{form.party2.title || ''}</td>
      </tr>
      <tr>
        <td class="row-label">Company</td>
        <td>{form.party1.company || ''}</td>
        <td>{form.party2.company || ''}</td>
      </tr>
      <tr>
        <td class="row-label">Notice Address</td>
        <td>{form.party1.address || ''}</td>
        <td>{form.party2.address || ''}</td>
      </tr>
      <tr>
        <td class="row-label">Date</td>
        <td></td>
        <td></td>
      </tr>
    </tbody>
  </table>

  <p class="cc-notice">
    Common Paper Mutual Non-Disclosure Agreement (Version 1.0) free to use under CC BY 4.0.
  </p>

  <!-- ══════════════════════════ STANDARD TERMS ══════════════════════════ -->
  <div class="page-break"></div>

  <h2>Standard Terms</h2>

  <ol class="standard-terms">
    <li>
      <strong>Introduction.</strong> This Mutual Non-Disclosure Agreement (which incorporates
      these Standard Terms and the Cover Page (defined below)) ("MNDA") allows each party
      ("Disclosing Party") to disclose or make available information in connection with the
      <em>{purpose}</em> which (1) the Disclosing Party identifies to the receiving party
      ("Receiving Party") as "confidential", "proprietary", or the like or (2) should be
      reasonably understood as confidential or proprietary due to its nature and the
      circumstances of its disclosure ("Confidential Information"). Each party's Confidential
      Information also includes the existence and status of the parties' discussions and
      information on the Cover Page. Confidential Information includes technical or business
      information, product designs or roadmaps, requirements, pricing, security and compliance
      documentation, technology, inventions and know-how. To use this MNDA, the parties must
      complete and sign a cover page incorporating these Standard Terms ("Cover Page"). Each
      party is identified on the Cover Page and capitalized terms have the meanings given herein
      or on the Cover Page.
    </li>

    <li>
      <strong>Use and Protection of Confidential Information.</strong> The Receiving Party
      shall: (a) use Confidential Information solely for the <em>{purpose}</em>; (b) not
      disclose Confidential Information to third parties without the Disclosing Party's prior
      written approval, except that the Receiving Party may disclose Confidential Information
      to its employees, agents, advisors, contractors and other representatives having a
      reasonable need to know for the <em>{purpose}</em>, provided these representatives are
      bound by confidentiality obligations no less protective of the Disclosing Party than the
      applicable terms in this MNDA and the Receiving Party remains responsible for their
      compliance with this MNDA; and (c) protect Confidential Information using at least the
      same protections the Receiving Party uses for its own similar information but no less
      than a reasonable standard of care.
    </li>

    <li>
      <strong>Exceptions.</strong> The Receiving Party's obligations in this MNDA do not apply
      to information that it can demonstrate: (a) is or becomes publicly available through no
      fault of the Receiving Party; (b) it rightfully knew or possessed prior to receipt from
      the Disclosing Party without confidentiality restrictions; (c) it rightfully obtained
      from a third party without confidentiality restrictions; or (d) it independently developed
      without using or referencing the Confidential Information.
    </li>

    <li>
      <strong>Disclosures Required by Law.</strong> The Receiving Party may disclose
      Confidential Information to the extent required by law, regulation or regulatory
      authority, subpoena or court order, provided (to the extent legally permitted) it
      provides the Disclosing Party reasonable advance notice of the required disclosure and
      reasonably cooperates, at the Disclosing Party's expense, with the Disclosing Party's
      efforts to obtain confidential treatment for the Confidential Information.
    </li>

    <li>
      <strong>Term and Termination.</strong> This MNDA commences on the
      <em>{effectiveDate}</em> and expires at the end of the <em>{mndaTerm}</em>. Either party
      may terminate this MNDA for any or no reason upon written notice to the other party. The
      Receiving Party's obligations relating to Confidential Information will survive for the
      <em>{confidentialityTerm}</em>, despite any expiration or termination of this MNDA.
    </li>

    <li>
      <strong>Return or Destruction of Confidential Information.</strong> Upon expiration or
      termination of this MNDA or upon the Disclosing Party's earlier request, the Receiving
      Party will: (a) cease using Confidential Information; (b) promptly after the Disclosing
      Party's written request, destroy all Confidential Information in the Receiving Party's
      possession or control or return it to the Disclosing Party; and (c) if requested by the
      Disclosing Party, confirm its compliance with these obligations in writing. As an
      exception to subsection (b), the Receiving Party may retain Confidential Information in
      accordance with its standard backup or record retention policies or as required by law,
      but the terms of this MNDA will continue to apply to the retained Confidential
      Information.
    </li>

    <li>
      <strong>Proprietary Rights.</strong> The Disclosing Party retains all of its intellectual
      property and other rights in its Confidential Information and its disclosure to the
      Receiving Party grants no license under such rights.
    </li>

    <li>
      <strong>Disclaimer.</strong> ALL CONFIDENTIAL INFORMATION IS PROVIDED "AS IS", WITH ALL
      FAULTS, AND WITHOUT WARRANTIES, INCLUDING THE IMPLIED WARRANTIES OF TITLE, MERCHANTABILITY
      AND FITNESS FOR A PARTICULAR PURPOSE.
    </li>

    <li>
      <strong>Governing Law and Jurisdiction.</strong> This MNDA and all matters relating
      hereto are governed by, and construed in accordance with, the laws of the State of
      <em>{governingLaw}</em>, without regard to the conflict of laws provisions of such
      <em>{governingLaw}</em>. Any legal suit, action, or proceeding relating to this MNDA
      must be instituted in the federal or state courts located in <em>{jurisdiction}</em>.
      Each party irrevocably submits to the exclusive jurisdiction of such
      <em>{jurisdiction}</em> in any such suit, action, or proceeding.
    </li>

    <li>
      <strong>Equitable Relief.</strong> A breach of this MNDA may cause irreparable harm for
      which monetary damages are an insufficient remedy. Upon a breach of this MNDA, the
      Disclosing Party is entitled to seek appropriate equitable relief, including an
      injunction, in addition to its other remedies.
    </li>

    <li>
      <strong>General.</strong> Neither party has an obligation under this MNDA to disclose
      Confidential Information to the other or proceed with any proposed transaction. Neither
      party may assign this MNDA without the prior written consent of the other party, except
      that either party may assign this MNDA in connection with a merger, reorganization,
      acquisition or other transfer of all or substantially all its assets or voting securities.
      Any assignment in violation of this Section is null and void. This MNDA will bind and
      inure to the benefit of each party's permitted successors and assigns. Waivers must be
      signed by the waiving party's authorized representative and cannot be implied from
      conduct. If any provision of this MNDA is held unenforceable, it will be limited to the
      minimum extent necessary so the rest of this MNDA remains in effect. This MNDA
      (including the Cover Page) constitutes the entire agreement of the parties with respect
      to its subject matter, and supersedes all prior and contemporaneous understandings,
      agreements, representations, and warranties, whether written or oral, regarding such
      subject matter. This MNDA may only be amended, modified, waived, or supplemented by an
      agreement in writing signed by both parties. Notices, requests and approvals under this
      MNDA must be sent in writing to the email or postal addresses on the Cover Page and are
      deemed delivered on receipt. This MNDA may be executed in counterparts, including
      electronic copies, each of which is deemed an original and which together form the same
      agreement.
    </li>
  </ol>

  <p class="cc-notice">
    Common Paper Mutual Non-Disclosure Agreement Version 1.0 free to use under CC BY 4.0.
  </p>
</div>

<style>
  .document {
    font-family: 'Georgia', 'Times New Roman', serif;
    font-size: 11pt;
    line-height: 1.6;
    color: #1a1a1a;
    padding: 48px 56px;
    background: #ffffff;
    max-width: 800px;
    margin: 0 auto;
  }

  h1 {
    font-size: 18pt;
    font-weight: bold;
    text-align: center;
    margin-bottom: 24px;
  }

  h2 {
    font-size: 13pt;
    font-weight: bold;
    margin-top: 28px;
    margin-bottom: 8px;
    border-bottom: 1px solid #333;
    padding-bottom: 4px;
  }

  h3 {
    font-size: 11pt;
    font-weight: bold;
    margin-top: 20px;
    margin-bottom: 4px;
  }

  p {
    margin: 6px 0;
  }

  .intro {
    margin-bottom: 16px;
  }

  .cover-field {
    margin-bottom: 12px;
  }

  .field-label {
    font-size: 9pt;
    font-style: italic;
    color: #555;
    margin: 0 0 4px 0;
  }

  .checkbox {
    font-size: 13pt;
    vertical-align: middle;
  }

  .checkbox.checked {
    color: #000;
  }

  .signing-intro {
    margin-top: 20px;
    margin-bottom: 12px;
    font-style: italic;
  }

  .signature-table {
    width: 100%;
    border-collapse: collapse;
    margin: 12px 0 20px;
    font-size: 10pt;
  }

  .signature-table th,
  .signature-table td {
    border: 1px solid #333;
    padding: 8px 10px;
    vertical-align: top;
  }

  .signature-table th {
    background-color: #f0f0f0;
    font-weight: bold;
    text-align: center;
  }

  .signature-table .row-label {
    font-weight: bold;
    background-color: #f8f8f8;
    width: 140px;
  }

  .signature-table tr:first-child td:not(.row-label),
  .signature-table tr:last-child td:not(.row-label) {
    height: 40px;
  }

  .cc-notice {
    font-size: 9pt;
    color: #666;
    margin-top: 16px;
    font-style: italic;
  }

  .page-break {
    height: 1px;
    border-top: 2px dashed #ccc;
    margin: 32px 0;
  }

  .standard-terms {
    padding-left: 24px;
    margin-top: 12px;
  }

  .standard-terms li {
    margin-bottom: 14px;
    text-align: justify;
  }

  em {
    font-style: normal;
    font-weight: bold;
    text-decoration: underline;
  }
</style>
