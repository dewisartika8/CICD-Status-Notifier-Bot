// API Testing Script - untuk testing koneksi dan endpoint
import { dashboardApi, projectApi, healthApi, webhookApi, telegramApi } from './src/services/api';

const apiTests = {
  // Test Health API
  async testHealth() {
    console.log('🔍 Testing Health API...');
    try {
      const response = await healthApi.checkHealth();
      console.log('✅ Health API:', response.data);
      return true;
    } catch (error: any) {
      console.log('❌ Health API Error:', error.message);
      return false;
    }
  },

  // Test Dashboard API
  async testDashboard() {
    console.log('🔍 Testing Dashboard API...');
    try {
      const response = await dashboardApi.getOverview();
      console.log('✅ Dashboard API:', response.data);
      return true;
    } catch (error: any) {
      console.log('❌ Dashboard API Error:', error.message);
      return false;
    }
  },

  // Test Projects API
  async testProjects() {
    console.log('🔍 Testing Projects API...');
    try {
      const response = await projectApi.getProjects();
      console.log('✅ Projects API:', response.data);
      return true;
    } catch (error: any) {
      console.log('❌ Projects API Error:', error.message);
      return false;
    }
  },

  // Test Telegram API
  async testTelegram() {
    console.log('🔍 Testing Telegram API...');
    try {
      const response = await telegramApi.getSubscriptions();
      console.log('✅ Telegram API:', response.data);
      return true;
    } catch (error: any) {
      console.log('❌ Telegram API Error:', error.message);
      return false;
    }
  },

  // Run all tests
  async runAllTests() {
    console.log('🚀 Starting API Integration Tests...\n');
    
    const results = {
      health: await this.testHealth(),
      dashboard: await this.testDashboard(),
      projects: await this.testProjects(),
      telegram: await this.testTelegram()
    };

    console.log('\n📊 Test Results Summary:');
    console.log('├─ Health API:', results.health ? '✅ PASS' : '❌ FAIL');
    console.log('├─ Dashboard API:', results.dashboard ? '✅ PASS' : '❌ FAIL');
    console.log('├─ Projects API:', results.projects ? '✅ PASS' : '❌ FAIL');
    console.log('└─ Telegram API:', results.telegram ? '✅ PASS' : '❌ FAIL');

    const passCount = Object.values(results).filter(Boolean).length;
    const totalTests = Object.keys(results).length;
    
    console.log(`\n🎯 Overall: ${passCount}/${totalTests} tests passed`);
    
    if (passCount === totalTests) {
      console.log('🎉 All tests passed! API integration successful!');
    } else {
      console.log('⚠️ Some tests failed. Check backend server status.');
    }

    return results;
  }
};

// Tambahkan ke window untuk akses dari browser console
if (typeof window !== 'undefined') {
  (window as any).apiTests = apiTests;
}

export default apiTests;
