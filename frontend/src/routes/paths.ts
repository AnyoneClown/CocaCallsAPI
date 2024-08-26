export const rootPaths = {
  homeRoot: '/',
  authRoot: 'authentication',
  errorRoot: 'error',
  loginRoot: '/authentication/login',
  profileRoot: 'profile',
};

export default {
  home: `/${rootPaths.homeRoot}`,
  login: `/${rootPaths.authRoot}/login`,
  signup: `/${rootPaths.authRoot}/sign-up`,
  oauthCallback: `/${rootPaths.authRoot}/callback`,
  404: `/${rootPaths.errorRoot}/404`,
  profile: `/${rootPaths.profileRoot}`,
};