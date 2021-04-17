import * as dotenv from "dotenv";

dotenv.config();

/**
 * Checks if a required environment variable is present and returns it if found.
 * If defaultValue is undefined an error will be thrown.
 *
 * @param key - The environment variable key.
 * @param defaultValue - A backup default value if the environment variable is not set.
 */
export function checkEnv(key: string, defaultValue?: string): string {
  if (key === "") {
    if (defaultValue !== undefined) {
      return defaultValue;
    }

    throw Error("Error: Failed to fetch environment variable: key cannot be empty");
  }

  const value = process.env[key];
  if (value === undefined || value === "") {
    if (defaultValue !== undefined) {
      return defaultValue;
    }

    throw Error(`Error: Failed to fetch environment variable for ${key}, value is ${value}`);
  }

  return value;
}