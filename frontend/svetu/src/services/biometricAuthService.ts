import { configManager } from '@/config';

interface BiometricAuthResult {
  success: boolean;
  method?: 'fingerprint' | 'face' | 'pin' | 'pattern';
  error?: string;
  timestamp: number;
}

interface BiometricCapabilities {
  available: boolean;
  fingerprint: boolean;
  faceId: boolean;
  webAuthn: boolean;
  touchId: boolean;
}

class BiometricAuthService {
  private capabilities: BiometricCapabilities | null = null;
  private publicKeyCredential: PublicKeyCredential | null = null;

  /**
   * Check device biometric capabilities
   */
  async checkCapabilities(): Promise<BiometricCapabilities> {
    if (this.capabilities) {
      return this.capabilities;
    }

    const capabilities: BiometricCapabilities = {
      available: false,
      fingerprint: false,
      faceId: false,
      webAuthn: false,
      touchId: false,
    };

    // Check WebAuthn support
    if (window.PublicKeyCredential) {
      capabilities.webAuthn = true;
      capabilities.available = true;

      // Check for platform authenticator (biometric)
      const available =
        await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
      if (available) {
        // Detect type based on user agent
        const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
        const isMac = /Macintosh/.test(navigator.userAgent);
        const isAndroid = /Android/.test(navigator.userAgent);
        const isWindows = /Windows/.test(navigator.userAgent);

        if (isIOS || isMac) {
          capabilities.touchId = true;
          capabilities.faceId =
            isIOS && !navigator.userAgent.includes('iPhone 8');
        } else if (isAndroid) {
          capabilities.fingerprint = true;
          capabilities.faceId = true; // Modern Android supports both
        } else if (isWindows) {
          capabilities.fingerprint = true;
          capabilities.faceId = true; // Windows Hello
        }
      }
    }

    // Check for legacy biometric APIs (for older devices)
    if (
      'credentials' in navigator &&
      'preventSilentAccess' in navigator.credentials
    ) {
      capabilities.available = true;
    }

    this.capabilities = capabilities;
    return capabilities;
  }

  /**
   * Register biometric authentication for a user
   */
  async register(
    userId: string,
    username: string
  ): Promise<BiometricAuthResult> {
    try {
      const capabilities = await this.checkCapabilities();
      if (!capabilities.webAuthn) {
        return {
          success: false,
          error: 'Biometric authentication not supported',
          timestamp: Date.now(),
        };
      }

      // Generate challenge from server
      const challenge = await this.getRegistrationChallenge(userId);

      // Create credential options
      const publicKeyCredentialCreationOptions: PublicKeyCredentialCreationOptions =
        {
          challenge: this.stringToArrayBuffer(challenge),
          rp: {
            name: 'Sve Tu Marketplace',
            id: window.location.hostname,
          },
          user: {
            id: this.stringToArrayBuffer(userId),
            name: username,
            displayName: username,
          },
          pubKeyCredParams: [
            { alg: -7, type: 'public-key' }, // ES256
            { alg: -257, type: 'public-key' }, // RS256
          ],
          authenticatorSelection: {
            authenticatorAttachment: 'platform',
            userVerification: 'required',
            requireResidentKey: false,
          },
          timeout: 60000,
          attestation: 'direct',
        };

      // Create credential
      const credential = (await navigator.credentials.create({
        publicKey: publicKeyCredentialCreationOptions,
      })) as PublicKeyCredential;

      if (!credential) {
        throw new Error('Failed to create credential');
      }

      // Save credential to server
      await this.saveCredential(userId, credential);

      // Store for later use
      this.publicKeyCredential = credential;

      return {
        success: true,
        method: this.detectBiometricMethod(),
        timestamp: Date.now(),
      };
    } catch (error) {
      console.error('Biometric registration error:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Registration failed',
        timestamp: Date.now(),
      };
    }
  }

  /**
   * Authenticate user with biometric
   */
  async authenticate(userId?: string): Promise<BiometricAuthResult> {
    try {
      const capabilities = await this.checkCapabilities();
      if (!capabilities.webAuthn) {
        return {
          success: false,
          error: 'Biometric authentication not supported',
          timestamp: Date.now(),
        };
      }

      // Get authentication challenge from server
      const { challenge, credentialIds } =
        await this.getAuthenticationChallenge(userId);

      // Create authentication options
      const publicKeyCredentialRequestOptions: PublicKeyCredentialRequestOptions =
        {
          challenge: this.stringToArrayBuffer(challenge),
          allowCredentials: credentialIds.map((id: string) => ({
            id: this.stringToArrayBuffer(id),
            type: 'public-key' as PublicKeyCredentialType,
            transports: ['internal' as AuthenticatorTransport],
          })),
          userVerification: 'required',
          timeout: 60000,
        };

      // Authenticate
      const credential = (await navigator.credentials.get({
        publicKey: publicKeyCredentialRequestOptions,
      })) as PublicKeyCredential;

      if (!credential) {
        throw new Error('Authentication failed');
      }

      // Verify with server
      const verified = await this.verifyCredential(credential);

      if (!verified) {
        throw new Error('Verification failed');
      }

      return {
        success: true,
        method: this.detectBiometricMethod(),
        timestamp: Date.now(),
      };
    } catch (error) {
      console.error('Biometric authentication error:', error);

      // Check if user cancelled
      if ((error as any).name === 'NotAllowedError') {
        return {
          success: false,
          error: 'Authentication cancelled',
          timestamp: Date.now(),
        };
      }

      return {
        success: false,
        error: error instanceof Error ? error.message : 'Authentication failed',
        timestamp: Date.now(),
      };
    }
  }

  /**
   * Enable quick unlock with PIN/Pattern as fallback
   */
  async enableQuickUnlock(pin?: string): Promise<boolean> {
    try {
      if (!pin || pin.length < 4) {
        throw new Error('PIN must be at least 4 digits');
      }

      // Hash PIN before storing
      const hashedPin = await this.hashPin(pin);

      // Store securely in IndexedDB or secure storage
      await this.secureStore('quickUnlockPin', hashedPin);

      return true;
    } catch (error) {
      console.error('Quick unlock setup error:', error);
      return false;
    }
  }

  /**
   * Verify quick unlock PIN
   */
  async verifyQuickUnlock(pin: string): Promise<boolean> {
    try {
      const storedHash = await this.secureRetrieve('quickUnlockPin');
      if (!storedHash) {
        return false;
      }

      const inputHash = await this.hashPin(pin);
      return storedHash === inputHash;
    } catch (error) {
      console.error('Quick unlock verification error:', error);
      return false;
    }
  }

  /**
   * Remove biometric authentication
   */
  async removeAuthentication(userId: string): Promise<boolean> {
    try {
      // Remove from server
      const apiUrl = configManager.get('api.url');
      await fetch(`${apiUrl}/api/v1/auth/biometric/remove`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userId }),
      });

      // Clear local storage
      await this.clearSecureStorage();

      this.publicKeyCredential = null;

      return true;
    } catch (error) {
      console.error('Remove authentication error:', error);
      return false;
    }
  }

  /**
   * Check if biometric is enrolled for user
   */
  async isEnrolled(userId: string): Promise<boolean> {
    try {
      const apiUrl = configManager.get('api.url');
      const response = await fetch(`${apiUrl}/api/v1/auth/biometric/enrolled/${userId}`);
      const { enrolled } = await response.json();
      return enrolled;
    } catch {
      return false;
    }
  }

  /**
   * Helper methods
   */
  private detectBiometricMethod(): 'fingerprint' | 'face' | 'pin' | 'pattern' {
    const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
    const isMac = /Macintosh/.test(navigator.userAgent);

    if (isIOS && !navigator.userAgent.includes('iPhone 8')) {
      return 'face';
    }

    if (isIOS || isMac) {
      return 'fingerprint'; // Touch ID
    }

    // Default to fingerprint for other platforms
    return 'fingerprint';
  }

  private async getRegistrationChallenge(userId: string): Promise<string> {
    const apiUrl = configManager.get('api.url');
    const response = await fetch(`${apiUrl}/api/v1/auth/biometric/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId, type: 'register' }),
    });
    const { challenge } = await response.json();
    return challenge;
  }

  private async getAuthenticationChallenge(userId?: string): Promise<{
    challenge: string;
    credentialIds: string[];
  }> {
    const apiUrl = configManager.get('api.url');
    const response = await fetch(`${apiUrl}/api/v1/auth/biometric/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId, type: 'authenticate' }),
    });
    return response.json();
  }

  private async saveCredential(
    userId: string,
    credential: PublicKeyCredential
  ): Promise<void> {
    const response = credential.response as AuthenticatorAttestationResponse;

    const apiUrl = configManager.get('api.url');
    await fetch(`${apiUrl}/api/v1/auth/biometric/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        userId,
        credentialId: this.arrayBufferToBase64(credential.rawId),
        publicKey: response.getPublicKey
          ? this.arrayBufferToBase64(response.getPublicKey()!)
          : '',
        attestation: this.arrayBufferToBase64(response.attestationObject),
      }),
    });
  }

  private async verifyCredential(
    credential: PublicKeyCredential
  ): Promise<boolean> {
    const response = credential.response as AuthenticatorAssertionResponse;

    const apiUrl = configManager.get('api.url');
    const verifyResponse = await fetch(`${apiUrl}/api/v1/auth/biometric/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        credentialId: this.arrayBufferToBase64(credential.rawId),
        authenticatorData: this.arrayBufferToBase64(response.authenticatorData),
        signature: this.arrayBufferToBase64(response.signature),
        userHandle: response.userHandle
          ? this.arrayBufferToBase64(response.userHandle)
          : null,
        clientDataJSON: this.arrayBufferToBase64(response.clientDataJSON),
      }),
    });

    const { verified } = await verifyResponse.json();
    return verified;
  }

  private stringToArrayBuffer(str: string): ArrayBuffer {
    const encoder = new TextEncoder();
    const encoded = encoder.encode(str);
    return encoded.buffer as ArrayBuffer;
  }

  private arrayBufferToBase64(buffer: ArrayBuffer): string {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  private async hashPin(pin: string): Promise<string> {
    const encoder = new TextEncoder();
    const data = encoder.encode(pin);
    const hash = await crypto.subtle.digest('SHA-256', data);
    return this.arrayBufferToBase64(hash);
  }

  private async secureStore(key: string, value: string): Promise<void> {
    // Use IndexedDB for secure storage
    const db = await this.openSecureDB();
    const tx = db.transaction('secure', 'readwrite');
    await tx.objectStore('secure').put({ key, value });
  }

  private async secureRetrieve(key: string): Promise<string | null> {
    const db = await this.openSecureDB();
    const tx = db.transaction('secure', 'readonly');
    const result = await tx.objectStore('secure').get(key);
    return (result as any)?.value || null;
  }

  private async clearSecureStorage(): Promise<void> {
    const db = await this.openSecureDB();
    const tx = db.transaction('secure', 'readwrite');
    await tx.objectStore('secure').clear();
  }

  private async openSecureDB(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open('BiometricSecure', 1);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result);

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains('secure')) {
          db.createObjectStore('secure', { keyPath: 'key' });
        }
      };
    });
  }
}

// Export singleton instance
export const biometricAuth = new BiometricAuthService();

// Export types
export type { BiometricAuthResult, BiometricCapabilities };
