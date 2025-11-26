import client from './client';

export const adminApi = {
    updateSystemSetting(key: string, value: string) {
        return client.post('/admin/settings/collector', {
            [key]: value, // Assuming the API expects a key-value pair or specific structure
        });
    },

    // Specific method for collector frequency as per backend implementation
    updateCollectorFrequency(runsPerDay: number) {
        return client.post('/admin/settings/collector', {
            runs_per_day: runsPerDay
        });
    }
};
