import type { SponsorList } from '@/generated/types';

const sponsorList: SponsorList = {
  contact: 'https://bipki.net',
  sponsors: [
    {
      id: '123',
      link: 'https://bipki.net',
      name: 'Bipki Services',
      slots: ['login', 'dashboard'],
      logo: '/assets/bipki-icon.jpg',
      text: { en: 'Bipki as always, remains unexplained' },
      title: { en: 'Bipki AI Suitcase Services' },
      until: '2099-01-01T00:00:00Z',
    },
  ],
};

export function useSponsorsQuery() {
  return { data: sponsorList, fetched: true };
}
