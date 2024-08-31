import { GridRowsProp } from '@mui/x-data-grid';
import { getToken } from 'api/auth';

export let rows: GridRowsProp = [];

export async function fetchCustomerData() {
  try {
    const token = getToken();

    const response = await fetch(`http://localhost:8080/api/users/`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      }
    });

    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    const responseData = await response.json();
    const data = responseData.data;
    rows = data.map((user: any) => ({
      id: user.ID,
      email: user.Email,
      provider: user.Provider,
      'verified-email': user.VerifiedEmail,
      'is-admin': user.IsAdmin,
      'created-at': user.CreatedAt,
      'updated-at': user.UpdatedAt,
      'deleted-at': user.DeletedAt,
      'subscription-id': user.Subscription?.ID ?? null,
      'subscription-start': user.Subscription?.StartDate ?? null,
      'subscription-end': user.Subscription?.EndDate ?? null,
    }));
  } catch (error) {
    console.error('Failed to fetch customer data:', error);
  }
}

fetchCustomerData();