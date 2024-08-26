interface UserMenuItem {
  id: number;
  path: string;
  title: string;
  icon: string;
  color?: string;
}

const userMenuItems: UserMenuItem[] = [
  {
    id: 1,
    path: '/profile',
    title: 'View Profile',
    icon: 'mingcute:user-2-fill',
    color: 'text.primary',
  },
  {
    id: 2,
    path: '!#',
    title: 'Notifications',
    icon: 'ion:notifications',
    color: 'text.primary',
  },
  {
    id: 3,
    path: '!#',
    title: 'Help Center',
    icon: 'material-symbols:live-help',
    color: 'text.primary',
  },
  {
    id: 4,
    path: '!#',
    title: 'Logout',
    icon: 'material-symbols:logout',
    color: 'error.main',
  },
];

export default userMenuItems;
