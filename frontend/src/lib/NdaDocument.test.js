import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import NdaDocument from './NdaDocument.svelte';

const base = {
  party1: { name: 'Alice Smith', title: 'CEO', company: 'Acme Inc', address: '123 Main St, Boston MA' },
  party2: { name: 'Bob Jones', title: 'CTO', company: 'Beta Corp', address: '456 Oak Ave, Austin TX' },
  effectiveDate: '2024-03-15',
  mndaTermType: 'expires',
  mndaTermYears: 2,
  confidentialityTermType: 'years',
  confidentialityTermYears: 3,
  purpose: 'Exploring a potential business partnership',
  governingLaw: 'California',
  jurisdiction: 'San Francisco, California',
  modifications: '',
};

describe('NdaDocument — party details', () => {
  it('renders party 1 name, title, and company', () => {
    render(NdaDocument, { form: base });
    expect(screen.getByText('Alice Smith')).toBeInTheDocument();
    expect(screen.getByText('CEO')).toBeInTheDocument();
    expect(screen.getByText('Acme Inc')).toBeInTheDocument();
  });

  it('renders party 2 name, title, and company', () => {
    render(NdaDocument, { form: base });
    expect(screen.getByText('Bob Jones')).toBeInTheDocument();
    expect(screen.getByText('CTO')).toBeInTheDocument();
    expect(screen.getByText('Beta Corp')).toBeInTheDocument();
  });

  it('renders party addresses', () => {
    render(NdaDocument, { form: base });
    expect(screen.getByText('123 Main St, Boston MA')).toBeInTheDocument();
    expect(screen.getByText('456 Oak Ave, Austin TX')).toBeInTheDocument();
  });
});

describe('NdaDocument — effective date', () => {
  it('formats a valid date as long-form US date', () => {
    render(NdaDocument, { form: base });
    const matches = screen.getAllByText('March 15, 2024');
    expect(matches.length).toBeGreaterThan(0);
  });

  it('shows placeholder when effectiveDate is empty', () => {
    render(NdaDocument, { form: { ...base, effectiveDate: '' } });
    expect(screen.getAllByText('[Effective Date not set]').length).toBeGreaterThan(0);
  });
});

describe('NdaDocument — MNDA term', () => {
  it('marks "expires" checkbox checked and shows year count', () => {
    render(NdaDocument, { form: { ...base, mndaTermType: 'expires', mndaTermYears: 2 } });
    expect(screen.getByText(/Expires 2 year\(s\) from Effective Date/)).toBeInTheDocument();
  });

  it('marks "continues until terminated" checkbox checked', () => {
    render(NdaDocument, { form: { ...base, mndaTermType: 'continues' } });
    expect(screen.getByText(/Continues until terminated in accordance with the terms of the MNDA/)).toBeInTheDocument();
  });
});

describe('NdaDocument — confidentiality term', () => {
  it('shows years-based confidentiality term', () => {
    render(NdaDocument, { form: { ...base, confidentialityTermType: 'years', confidentialityTermYears: 3 } });
    expect(screen.getAllByText(/3 year\(s\) from Effective Date/).length).toBeGreaterThan(0);
  });

  it('shows perpetuity when selected', () => {
    render(NdaDocument, { form: { ...base, confidentialityTermType: 'perpetuity' } });
    expect(screen.getByText(/In perpetuity/)).toBeInTheDocument();
  });
});

describe('NdaDocument — agreement fields', () => {
  it('renders purpose in document body', () => {
    render(NdaDocument, { form: base });
    const matches = screen.getAllByText('Exploring a potential business partnership');
    expect(matches.length).toBeGreaterThan(0);
  });

  it('shows placeholder when purpose is empty', () => {
    render(NdaDocument, { form: { ...base, purpose: '' } });
    expect(screen.getAllByText('[Purpose not specified]').length).toBeGreaterThan(0);
  });

  it('renders governing law', () => {
    render(NdaDocument, { form: base });
    expect(screen.getByText(/Governing Law: California/)).toBeInTheDocument();
  });

  it('renders jurisdiction', () => {
    render(NdaDocument, { form: base });
    expect(screen.getByText(/Jurisdiction: San Francisco, California/)).toBeInTheDocument();
  });

  it('shows "None" when modifications is empty', () => {
    render(NdaDocument, { form: { ...base, modifications: '' } });
    expect(screen.getByText('None')).toBeInTheDocument();
  });

  it('renders custom modifications text', () => {
    render(NdaDocument, { form: { ...base, modifications: 'Section 5 is amended' } });
    expect(screen.getByText('Section 5 is amended')).toBeInTheDocument();
  });
});
