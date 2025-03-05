import React from 'react';
import { FaExclamationTriangle } from 'react-icons/fa';

interface ErrorMessageProps {
  message: string;
  className?: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, className = '' }) => {
  if (!message) return null;

  return (
    <div className={`text-red-600 flex items-center mt-1 text-sm ${className}`}>
      <FaExclamationTriangle className="mr-1" />
      <span>{message}</span>
    </div>
  );
};

export default ErrorMessage; 