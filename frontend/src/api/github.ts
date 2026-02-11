import client from './client';

export async function disconnectGitHub() {
  return await client.delete('/auth/github/disconnect');
}
