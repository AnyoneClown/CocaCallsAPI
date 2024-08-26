import React, { useEffect, useState, ReactElement } from 'react';
import { getToken, getUserIDFromToken } from 'api/auth';
import { Box, Grid, Paper, Typography, Avatar, Divider } from '@mui/material';
import { Email, Verified, AdminPanelSettings, CalendarToday } from '@mui/icons-material';

interface Subscription {
  ID: string;
  UserID: string;
  StartDate: string;
  EndDate: string;
}

interface User {
  ID: string;
  Email: string;
  Picture: string;
  Provider: string;
  VerifiedEmail: boolean;
  IsAdmin: boolean;
  CreatedAt: string;
  Subscription: Subscription;
}

const ProfilePage = (): ReactElement => {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const token = getToken();
        const userID = getUserIDFromToken(token);
        if (userID === null || userID === undefined) {
          throw new Error('Invalid token');
        }
        const response = await fetch(`http://localhost:8080/api/users/${userID}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) {
          throw new Error('Failed to fetch user data');
        }
        const result = await response.json();
        const data = result.data;
        if (data) {
          setUser(data);
        }
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    };

    fetchUserData();
  }, []);

  if (!user) {
    return <Typography>Loading...</Typography>;
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('uk-UA', { year: 'numeric', month: '2-digit', day: '2-digit' });
  };

  return (
    <Box sx={{ p: 3, backgroundColor: '#1e1e2d', color: 'white', minHeight: '100vh' }}>
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Paper elevation={3} sx={{ p: 3, backgroundColor: '#2a2a3c', height: '100%' }}>
            <Box display="flex" flexDirection="column" alignItems="center">
              <Avatar src={user.Picture} alt="User Profile" sx={{ width: 120, height: 120, mb: 2 }} />
              <Typography variant="h5" color="common.white" sx={{ wordBreak: 'break-word', textAlign: 'center' }}>{user.Email}</Typography>
              <Typography variant="body2" color="textSecondary">{user.Provider}</Typography>
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={12} md={8}>
          <Paper elevation={3} sx={{ p: 3, backgroundColor: '#2a2a3c', height: '100%' }}>
            <Typography variant="h6" color="common.white"gutterBottom>User Information</Typography>
            <Box sx={{ mt: 2 }}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Box display="flex" alignItems="center" mb={2}>
                    <Email sx={{ mr: 1, color: 'white' }} />
                    <Box>
                      <Typography variant="body2" color="grey.400">Email</Typography>
                      <Typography variant="body1">{user.Email}</Typography>
                    </Box>
                  </Box>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Box display="flex" alignItems="center" mb={2}>
                    <Verified color={user.VerifiedEmail ? "primary" : "disabled"} sx={{ mr: 1 }} />
                    <Box>
                      <Typography variant="body2" color="grey.400">Verified Email</Typography>
                      <Typography variant="body1">{user.VerifiedEmail ? 'Yes' : 'No'}</Typography>
                    </Box>
                  </Box>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Box display="flex" alignItems="center" mb={2}>
                    <AdminPanelSettings color={user.IsAdmin ? "primary" : "disabled"} sx={{ mr: 1 }} />
                    <Box>
                      <Typography variant="body2" color="grey.400">Admin Status</Typography>
                      <Typography variant="body1">{user.IsAdmin ? 'Yes' : 'No'}</Typography>
                    </Box>
                  </Box>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Box display="flex" alignItems="center" mb={2}>
                    <CalendarToday sx={{ mr: 1, color: 'white' }} />
                    <Box>
                      <Typography variant="body2" color="grey.400">Created At</Typography>
                      <Typography variant="body1">{formatDate(user.CreatedAt)}</Typography>
                    </Box>
                  </Box>
                </Grid>
              </Grid>
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper elevation={3} sx={{ p: 3, backgroundColor: '#2a2a3c' }}>
            <Typography variant="h6" color="common.white" gutterBottom>Subscription Details</Typography>
            <Divider sx={{ my: 2, backgroundColor: 'rgba(255,255,255,0.1)' }} />
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="grey.400">Start Date</Typography>
                <Typography variant="body1">{formatDate(user.Subscription.StartDate)}</Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="grey.400">End Date</Typography>
                <Typography variant="body1">{formatDate(user.Subscription.EndDate)}</Typography>
              </Grid>
            </Grid>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default ProfilePage;``