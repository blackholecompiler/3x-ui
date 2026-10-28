import { useQuery } from '@tanstack/react-query';

import { HttpUtil } from '@/utils';
import { keys } from '@/api/queryKeys';
import type { SponsorList } from '@/generated/types';

const EMPTY: SponsorList = { sponsors: [] };

async function fetchSponsors(): Promise<SponsorList> {
  const msg = await HttpUtil.get<SponsorList>('/sponsors', undefined, { silent: true });
  if (!msg?.success || !msg.obj) return EMPTY;
  return { contact: msg.obj.contact, sponsors: msg.obj.sponsors ?? [] };
}

export function useSponsorsQuery() {
  // const query = useQuery({
  //   queryKey: keys.sponsors(),
  //   queryFn: fetchSponsors,
  //   staleTime: 60 * 60 * 1000,
  //   retry: false,
  // });

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

  return { data: sponsorList, fetched: true };
}
