/*  useEffect(() => {
    const url = "ws://localhost:3001/rtc";

    const chat = new Chat();

    chat.connect(url, "token");
  }, []); */

/*  async function sendKey() {
    const params = {
      name: "RSA-PSS",
      modulusLength: 2048, // The size of the RSA key
      publicExponent: new Uint8Array([0x01, 0x00, 0x01]), // 65537 is a common choice
      hash: "SHA-256", // The hash algorithm to use for signing
    };

    const keyPair = await crypto.subtle.generateKey(params, true, [
      "sign",
      "verify",
    ]);

    const spki = await crypto.subtle.exportKey("spki", keyPair.publicKey);

    const binarySpki = Array.from(new Uint8Array(spki));
    const base64Spki = btoa(String.fromCharCode.apply(null, binarySpki));
    const publicKey = `-----BEGIN PUBLIC KEY-----\n${base64Spki}\n-----END PUBLIC KEY-----`;

    fetch("http://localhost:3001/auth/complete", {
      method: "POST",
      body: JSON.stringify({ publicKey }),
      headers: {
        "Content-Type": "application/json",
      },
    });
  } */
